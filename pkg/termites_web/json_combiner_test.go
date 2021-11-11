package termites_web

import "testing"

func TestCombiner(t *testing.T) {
	combiner := newCombiner()

	var actual []byte
	actual, _ = combiner.get()

	res0 := "{\"version\":0,\"fields\":{}}"
	if string(actual) != res0 {
		t.Errorf("incorrect result\nexpected: %s\ngot:      %s\n", res0, actual)
	}

	combiner.update(JsonPartialData{
		Key:  "key 1",
		Data: []byte("{\"data\": false}"),
	})

	actual, _ = combiner.get()
	res1 := "{\"version\":1,\"fields\":{\"key 1\":{\"data\":false}}}"
	if string(actual) != res1 {
		t.Errorf("incorrect result\nexpected: %s\ngot:      %s\n", res1, actual)
	}

	combiner.update(JsonPartialData{
		Key:  "key 2",
		Data: []byte("{\"data\": \"check\"}"),
	})

	actual, _ = combiner.get()
	res2 := "{\"version\":2,\"fields\":{\"key 1\":{\"data\":false},\"key 2\":{\"data\":\"check\"}}}"
	if string(actual) != res2 {
		t.Errorf("incorrect result\nexpected: %s\ngot:      %s\n", res2, actual)
	}

	combiner.update(JsonPartialData{
		Key:  "key 1",
		Data: []byte("{\"data\": true}"),
	})

	actual, _ = combiner.get()
	res3 := "{\"version\":3,\"fields\":{\"key 1\":{\"data\":true},\"key 2\":{\"data\":\"check\"}}}"
	if string(actual) != res3 {
		t.Errorf("incorrect result\nexpected: %s\ngot:      %s\n", res3, actual)
	}
}
