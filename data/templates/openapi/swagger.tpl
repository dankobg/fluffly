<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <meta name="description" content="SwaggerUI" />
    <title>Fluffly SwaggerUI</title>
    <link
      rel="stylesheet"
      href="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui.css"
    />
    <script
      src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-bundle.js"
      crossorigin
      async
      defer
    ></script>
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script>
      window.onload = () => {
        window.ui = SwaggerUIBundle({
          url: "{{ .SpecURL }}",
          dom_id: "#swagger-ui",
        });
      };
    </script>
  </body>
</html>
