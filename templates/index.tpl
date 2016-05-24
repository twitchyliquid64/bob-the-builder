
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

          #overlay {
            width: 100%;
            height: 100%;
            position: absolute;
            top: 0;
            left: 0;
            z-index: 99999;
          }

      </style>

  </head>

  <body ng-app="baseApp" class="">
    <div class="ui dimmer" id="overlay">
        <div class="ui indeterminate text loader">Loading build data</div>
    </div>

    <div class="ui grid">
      <div class="four wide column">
        {!{template "topnav" .}!}
      </div>
      <div class="eleven wide column">
        <div ng-view></div>
      </div>
      <div class="ui one wide column"/>
    </div>


    {!{template "tailcontent"}!}

  </body>
</html>
