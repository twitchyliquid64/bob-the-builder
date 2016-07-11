(function () {

  angular.module('baseApp')
    .service('eventService', ['$rootScope', function($rootScope) {
      var self = this;
      var location = window.location;
      self.connected = false;
      var ws = new WebSocket("ws://" + location.hostname+(location.port ? ':'+location.port: '') + "/ws/events");

      ws.onopen = function()
      {
        console.log("Events ws opened.");
        $rootScope.$apply(function(){
          self.connected = true;
          $rootScope.$broadcast('ws-events-connected');
        });
      };

      ws.onmessage = function (evt)
      {
        var received_msg = evt.data;
        $rootScope.$apply(function(){
          var msg = JSON.parse(evt.data);
          if (msg.Type != "SERVER-STATS")console.log(msg);
          if (msg.Type == "RUN-STARTED"){
            $rootScope.$broadcast('ws-events-run-started', {run: msg.Data, index: msg.Index});
            $rootScope.$broadcast('ws-event-run-start-' + msg.Index);
          } else if (msg.Type == "RUN-FINISHED"){
            $rootScope.$broadcast('ws-events-run-finished', {run: msg.Data, index: msg.Index});
            $rootScope.$broadcast('ws-event-run-finish-' + msg.Index);
          } else if (msg.Type == "PHASE-STARTED"){
            $rootScope.$broadcast('ws-events-phase-started', {run: msg.Data, index: msg.Index});
            $rootScope.$broadcast('ws-event-phase-started-' + msg.Index, {phase: msg.Data});
          } else if (msg.Type == "PHASE-FINISHED"){
            $rootScope.$broadcast('ws-events-phase-finished', {run: msg.Data, index: msg.Index});
            $rootScope.$broadcast('ws-event-phase-finished-' + msg.Index, {phase: msg.Data});
          } else if (msg.Type == "PHASE-DATA"){
            $rootScope.$broadcast('ws-events-phase-data', {content: msg.Data.content, phase: msg.Data.phase, index: msg.Index});
            $rootScope.$broadcast('ws-event-phase-data-' + msg.Index, {content: msg.Data.content, phase: msg.Data.phase, index: msg.Index});
          } else if (msg.Type == "RELOAD-QUEUED"){
            $rootScope.$broadcast('ws-events-reload-queued', {});
          } else if (msg.Type == "DEF-REFRESH"){
            $rootScope.$broadcast('ws-events-reload-started', {});
          } else if (msg.Type == "DEF-REFRESH-COMPLETED"){
            $rootScope.$broadcast('ws-events-reload-finished', {});
          } else if (msg.Type == "SERVER-STATS"){
            $rootScope.$broadcast('ws-server-stats', msg.Data);
          }

        });
      };

      ws.onclose = function()
      {
        console.log("event ws closed.");
        $rootScope.$apply(function(){
          self.connected = false;
          $rootScope.$broadcast('ws-events-closed');
        });
      };

      ws.onerror = function(ev){
        console.log("event ws error:", ev);
      }
    }]);
})();
