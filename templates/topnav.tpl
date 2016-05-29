
    <div class="ui fluid large vertical menu" style="height:98vh; border-bottom: none;">
      <div class="ui header item">
        <i class="server big icon"></i>
        {!{.Config.Name}!}
      </div>

      <a class="active item">
        Dashboard
      </a>

      <div class="item">
        <div class="header">Build Definitions</div>
        <div class="menu">
            <a class="item" ng-repeat="definition in dataService.getDefinitions()" ng-click="navBuild($index)">
              {{definition.name}}
            </a>
        </div>
      </div>

      <a class="item">
        System Log
      </a>
      <div class="item">
        {{!dataService.connectionEstablished && !dataService.connectionLost? "Connecting ..." : ""}}
        {{dataService.connectionLost && !dataService.error ? "Connection Lost" : ""}}
        {{dataService.connectionEstablished && !dataService.connectionLost ? "Connected" : ""}}
        {{dataService.error ? "Service Error" : ""}}
        <i class="icon" ng-class="{yellow: dataService.error, red: (dataService.connectionLost && !dataService.error), checkmark: (dataService.connectionEstablished && !dataService.connectionLost), sign: dataService.connectionLost, spinner: !dataService.connectionEstablished && !dataService.connectionLost, loading: !dataService.connectionEstablished && !dataService.connectionLost, warning: dataService.connectionLost}"></i>
      </div>

    </div>
