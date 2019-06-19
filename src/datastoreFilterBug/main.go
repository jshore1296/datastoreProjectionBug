package main

// modified from HelloWorld example at https://github.com/GoogleCloudPlatform/golang-samples/blob/master/appengine/go11x/helloworld/helloworld.go
import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/datastore"
)

func main() {
	http.HandleFunc("/", indexHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

type MyStruct struct {
	Name    string
	Created time.Time
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	ctx := r.Context()
	project := os.Getenv("GOOGLE_CLOUD_PROJECT")
	dsClient, err := datastore.NewClient(ctx, project)
	if err != nil {
		_, _ = fmt.Fprint(w, "Creating Client", err)
		return
	}

	for i := 0; i < 5; i++ {
		k := datastore.IncompleteKey("MyStruct", nil)
		k, err := dsClient.Put(ctx, k, &MyStruct{
			Name:    fmt.Sprintf("Struct-%d", i),
			Created: time.Now().Add(time.Duration(-1 * i) * time.Hour),
		})
		if err != nil {
			_, _ = fmt.Fprint(w, "Creating MyStruct", err)
			return
		}
		defer func() {
			_ = dsClient.Delete(ctx, k)
		}()
	}
	k := datastore.IncompleteKey("MyStruct", nil)
	k, err = dsClient.Put(ctx, k, &MyStruct{
		Name:    fmt.Sprintf("Struct-%d", 10),
		Created: time.Time{},
	})
	if err != nil {
		_, _ = fmt.Fprint(w, "Creating MyStruct", err)
		return
	}
	defer func() {
		_ = dsClient.Delete(ctx, k)
	}()

	res1 := make([]MyStruct, 0)
	badQuery := datastore.NewQuery("MyStruct").Filter("Created > ", time.Now().Add(-6 * time.Hour))
	keys, err := dsClient.GetAll(ctx, badQuery, &res1)
	if err != nil {
		_, _ = fmt.Fprint(w, "BadQuery", err)
		return
	}
	_, _ = fmt.Fprintf(w, "BadQuery:\nKeys: %+v\nValues: %+v\n", keys, res1)
}
