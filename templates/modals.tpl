

<div class="ui modal" id="miniDocumentationModal">
  <div class="header">
    Documentation
  </div>

  <div class="content" style="min-height: 450px;">
    <div class="ui fuild accordion" id="documentation-defedit" style="background-color: white;">
      <div class="title">
        <i class="dropdown icon"></i>
        What are the important fields?
      </div>
      <div class="content">
        <p>A dog is a type of domesticated animal. Known for its loyalty and faithfulness, it can be found as a welcome guest in many households across the world.</p>
      </div>
      <div class="title">
        <i class="dropdown icon"></i>
        What kinds of steps can I use?
      </div>
      <div class="content">
        <p>There are many breeds of dogs. Each breed varies in size and temperament. Owners often select a breed of dog that they find to be compatible with their own lifestyle and desires from a companion.</p>
      </div>
      <div class="title">
        <i class="dropdown icon"></i>
        What kind of parameters can I use?
      </div>
      <div class="content">
        <p>Three common ways for a prospective owner to acquire a dog is from pet shops, private owners, or shelters.</p>
        <p>A pet shop may be the most convenient way to buy a dog. Buying a dog from a private owner allows you to assess the pedigree and upbringing of your dog before choosing to take it home. Lastly, finding your dog from a shelter, helps give a good home to a dog who may not find one so readily.</p>
      </div>
      <div class="title">
        <i class="dropdown icon"></i>
        Tempate Snippets
      </div>
      <div class="content">
        <p>A dog is a type of domesticated animal. Known for its loyalty and faithfulness, it can be found as a welcome guest in many households across the world.</p>
      </div>
    </div>
  </div>

  <div class="actions" style="text-align: initial; display: flex;">
    <div style="flex: 1; text-align: right;">
      <div class="ui button cancel">Close</div>
    </div>
  </div>
</div>









<div class="ui modal" id="runOptionsModal" ng-controller="runOptionsController">
  <div class="header">
    Run Options
  </div>

  <div class="content">

    <h4 class="ui dividing header">
      <i class="options icon"></i>
      <div class="content">
        Run parameters <i class="notched circle loading icon" ng-if="loadingParams"></i>
      </div>
    </h4>
    <div class="ui basic segment form"  style="min-height: 50px;">
      <div ng-repeat="param in defObj.params">

        <div ng-if="param.type == 'text'" class="field" style="margin-top: 7px;">
          <label>{{param.label}}</label>
          <input type="text" name="{{param.varname}}" ng-model="param.default" placeholder="{{param.placeholder}}" id="{{'runopt-field-' + $index}}">
        </div>


        <div ng-if="param.type == 'check'" class="field" style="margin-top: 7px;">
          <div class="ui checkbox">
            <input type="checkbox" tabindex="0" ng-model="param.default" id="{{'runopt-field-' + $index}}">
            <label>{{param.label}}</label>
          </div>
        </div>

        <div ng-if="param.type == 'branchselect'" style="margin-top: 7px;">
          <div class="ui fluid search selection dropdown" id="{{'runopt-field-' + $index}}">
            <input type="hidden" name="{{'nm-runopt-field-' + $index}}">
            <i class="dropdown icon"></i>
            <div class="default text">{{param.label}}</div>
            <div class="menu">
            </div>
          </div>
        </div>

        <div ng-if="param.type == 'select'" style="margin-top: 7px;">
          <div class="ui fluid search selection dropdown" id="{{'runopt-field-' + $index}}">
            <input type="hidden" name="{{'nm-runopt-field-' + $index}}">
            <i class="dropdown icon"></i>
            <div class="default text">{{param.label}}</div>
            <div class="menu">
              <div class="item" data-value="{{value}}" ng-repeat="(display, value) in param.items">
                {{display}}
              </div>
            </div>
          </div>
        </div>

      </div>
      <p ng-if="defObj.params == null || defObj.params == undefined || defObj.params.length == 0">No run parameters exist in the build definition.</p>
    </div>



    <h4 class="ui dividing header">
      <i class="tag icon"></i>
      <div class="content">
        Tags
      </div>
    </h4>

    <div class="ui basic segment" style="min-height: 50px;">
      <div class="ui fluid multiple search selection dropdown" style="border: 0px;" id="runOptionsModal-tagsDropdown">
        <input name="tags" type="hidden">
        <div class="default text"><i class="tags icon"></i> Add tags...</div>
        <div class="menu">
          <div class="item" data-value="production">
            <div class="ui red empty circular label"></div>
            Production
          </div>
          <div class="item" data-value="staging">
            <div class="ui yellow empty circular label"></div>
            Staging
          </div>
          <div class="item" data-value="test">
            <div class="ui green empty circular label"></div>
            Test
          </div>
          <div class="item" data-value="backup">
            <div class="ui purple empty circular label"></div>
            Backup
          </div>
        </div>
      </div>
    </div>

    <h4 class="ui dividing header">
      <i class="adjust icon"></i>
      <div class="content">
        Auxillary
      </div>
    </h4>

    <div class="ui basic segment" style="min-height: 50px;">
      <div class="ui checkbox">
        <input type="checkbox" name="disphys" id="runOptionsModal-disablephys">
        <label>Disable physical indicators</label>
      </div>
      <div class="ui checkbox" style="margin-left: 3em;">
        <input type="checkbox" name="preventpostbuild">
        <label>Prevent post-build triggers</label>
      </div>
    </div>
  </div>

  <div class="actions" style="text-align: initial; display: flex;">

    <div style="flex: 1; text-align: left;">
      <div class="ui labeled input">
        <div class="ui label">
          version
        </div>
        <input type="text" name="version" placeholder="0.0.1" ng-model="defObj['last-version']" style="width: 78px;">
      </div>
    </div>


    <div style="flex: 1; text-align: right;">
      <div class="ui button cancel">Cancel</div>
      <div class="ui button primary icon ok">
        <i class="icon rocket"></i>
        Run
      </div>
    </div>
  </div>
</div>
