<style>
.definition-refresh {
  cursor: pointer;
  cursor: hand;
}
</style>


    <div class="ui fluid large vertical menu" style="height:98vh; border-bottom: none;">
      <div class="ui header item">
        <i class="definition-refresh refresh big icon" ng-click="dataService.requestDefinitionsReload()"></i>
        {!{.Config.Name}!}
      </div>

      <a class="item" ng-class="{active: currentlyDash}" ng-click="navDashboard()">
        Dashboard
      </a>

      <div class="item">
        <div class="header">Definitions</div>
        <div class="menu">
            <a class="item" ng-repeat="definition in dataService.getDefinitions()" ng-click="navBuild($index)" ng-class="{active: !currentlyDash && (currentIndex == $index)}">
              <i class="icon" ng-class="definition.icon"></i>&nbsp;&nbsp;&nbsp;{{definition.name}}
            </a>
        </div>
      </div>

      <a class="item" ng-class="{active: currentlyDocumentation}" ng-click="documentation()">
        Documentation
      </a>
      <a class="item" ng-class="{active: currentlyBrowser}" ng-click="browser()">
        MDS Browser
      </a>
      <a class="item" ng-class="{active: currentlyCron}" ng-click="cron()">
        Cron entries
      </a>


      <div class="item" id="connStatusText">
        {{!dataService.connectionEstablished && !dataService.connectionLost? "Connecting ..." : ""}}
        {{dataService.connectionLost && !dataService.error ? "Connection Lost" : ""}}
        {{dataService.connectionEstablished && !dataService.connectionLost ? "Connected" : ""}}
        {{dataService.error ? "Service Error" : ""}}
        <i class="icon" ng-class="{yellow: dataService.error, red: (dataService.connectionLost && !dataService.error), checkmark: (dataService.connectionEstablished && !dataService.connectionLost), sign: dataService.connectionLost, loading: !dataService.connectionEstablished && !dataService.connectionLost, notched: !dataService.connectionEstablished && !dataService.connectionLost, circle: !dataService.connectionEstablished && !dataService.connectionLost, warning: dataService.connectionLost}"></i>
      </div>
      <div class="item" ng-if="dataService.reloadQueued">
        Definitions reloading ... Please wait.
        <i class="icon yellow sign"></i>
      </div>

    </div>
