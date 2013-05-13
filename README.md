Open-Voice-Data-Project-Website
===============================

Update
------

This past weekend, Travis took it upon himself to clean up the PHP site (which was almost a static site anyways, with about 10 lines of PHP interpolated into HTML), and separate everything into a cleanly templated, Go-backed site.  The routing and template rendering code is in main.go, while a small example of consuming the API is in funcs.go.