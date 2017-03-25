
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

          #bottom-right-info{
            position: fixed;
            bottom: 5px;
            right: 7px;
            padding: 6px;
            text-align: right;
          }
      </style>

  </head>

  <body ng-app="baseApp" ng-controller="mainController">
    <div class="ui dimmer" id="overlay" ng-class="{active: dataService.isLoading()}">
        <div class="ui indeterminate text loader">{{dataService.loadingMessage}}</div>
    </div>

    <div class="ui grid">
      <div class="four wide column" id="sidebar">
        {!{template "topnav" .}!}
      </div>
      <div class="eleven wide column">
        <div ng-view></div>
      </div>
      <div class="ui one wide column"/>
    </div>

    {!{template "modals"}!}

    {!{template "tailcontent"}!}

    {!{if .Auth}!}
    <div id="bottom-right-info">
      <p>
        <b>User:</b> <i>{!{.Auth.User.Name}!}</i><br />
        <b>Session created:</b> <i>{!{.Auth.User.Created.Format "Jan _2 15:04:05"}!}</i><br />
        <b>Session expires:</b> <i>{!{.Auth.User.Expiry.Format "Jan _2 15:04:05"}!}</i>
        <br />Password {!{if .Auth.User.OTPUsed}!}+ <span style="color: #21ba45;">OTP</span> {!{end}!}authentication
      </p>
    </div>
    {!{end}!}
  </body>
</html>
