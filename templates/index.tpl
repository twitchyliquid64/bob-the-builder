
<!DOCTYPE html>
<html lang="en">

  <head>
      <title>{!{.Config.Name}!}</title>
      {!{template "headcontent"}!}

      <style type="text/css">
          /**
           * Hide when Angular is not yet loaded and initialized
           */
          [ng\:cloak], [ng-cloak], [data-ng-cloak], [x-ng-cloak], .ng-cloak, .x-ng-cloak {
            display: none !important;
          }
      </style>

  </head>

  <body id="example" class="started" ontouchstart="">



    {!{template "topnav" .}!}




    {!{template "tailcontent"}!}

  </body>
</html>
