<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Fluffly Stoplight</title>
    <link
      rel="stylesheet"
      href="https://unpkg.com/@stoplight/elements/styles.min.css"
    />
    <script
      src="https://unpkg.com/@stoplight/elements/web-components.min.js"
      async
      defer
    ></script>
  </head>
  <body>
    <elements-api
      apiDescriptionUrl="{{ .SpecURL }}"
      router="hash"
      layout="sidebar"
    />
  </body>
</html>
