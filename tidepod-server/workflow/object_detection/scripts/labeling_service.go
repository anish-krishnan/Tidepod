package script

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// GetLabelsForImage queries the python classification service
// to get all the labels for an image
func GetLabelsForImage(filename string) ([]string, error) {
	// Encode the data
	form := url.Values{}
	form.Add("filename", filename)

	resp, err := http.PostForm("http://localhost:2999/labelImage", form)

	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	defer resp.Body.Close()

	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	type Data struct {
		Labels []string
	}

	data := &Data{Labels: []string{}}
	err = json.Unmarshal(body, data)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return data.Labels, nil
}
