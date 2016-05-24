
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
          {!{range .Builder.Definitions}!}
            <a class="item">
              {!{.Name}!}
            </a>
          {!{end}!}
        </div>
      </div>

      <a class="item">
        System Log
      </a>
      <div class="item">
        Connected
      </div>

    </div>
