package parser_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nkanaev/yarr/src/parser"
)

func TestOpenGraph(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
		rw.WriteHeader(http.StatusOK)
		fmt.Fprint(rw, html)
	}))
	t.Cleanup(ts.Close)

	og, err := parser.ParseOpenGraph(ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if og["image"] != "https://nadia.moe/images/nadia.moe-banner.gif" {
		t.Fatalf("Unexpected image %q", og["image"])
	}
}

var html string = `
<!DOCTYPE html>
<html lang="en">

<head>
  <meta name="title" property="og:title" content="Nadia.moe">
  <meta name="description" property="og:description" content="Nadia's links and contact information.">
  <meta name="author" content="Nadia Santalla">
  <meta name="image" property="og:image" content="https://nadia.moe/images/nadia.moe-banner.gif">
  <meta property="og:image:type" content="image/gif">
  <meta property="og:image:width" content="1320">
  <meta property="og:image:height" content="690">
  <meta property="og:url" content="https://nadia.moe">
  <meta name="fediverse:creator" content="@nadia@raru.re">
</head>

<body class="monospace">
</body>

</html>
`
