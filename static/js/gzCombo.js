/**
 * gzCombo Version 1.0
 * Created by Greg Zhang on 2016/6/28.
 */
;(function(window, $){
    var gzCombo,
        defaults = {
            "width" : "100%",
            "dataSource" : null, // String or Array
            "requestData" : null, //Depends on dataSource
            "placeholder" : "",
            "displayMode" : "block",
            "emptyValue" : "",
            "isAsync" : false,
            "onValueChange" : null,
            "keyMap" : {
                "ENTER" : 13,
                "UP" : 38,
                "DOWN" : 40
            }
        },
        keyIndex = -1,
        isLoaded = false;

    function _processRenderData(dataSource){
        for(var i = 0; i < dataSource.length; i++){
            var comboDataObj = dataSource[i];
            $('<li class="gzCombo-drop-list-item">'+ comboDataObj.labelName +'</li>')
                .data({
                    'labelValue' : comboDataObj.labelValue,
                    'labelName' : comboDataObj.labelName,
                    'gzIndex' : i
                }).attr('data-gz-index', i).appendTo(this.dropList);
        }

    }

    function _processRenderedStyle(){
        this.dropList.css({
            "width" : this.element.outerWidth(),
            "left" : "0px",
            "bottom" : -(this.dropList.outerHeight() + 2) + "px"
        }).show();
    }

    gzCombo = function(options){
        var gzComboPlugin = this;

        $.extend(defaults, options || {});
        this.getText = function () {
            return this.inputEl.val();
        };

        this.getValue = function(){
            return this.inputEl.data('labelValue') ? this.inputEl.data('labelValue') : defaults.emptyValue;
        };

        this.getIndex = function($item){
            var index = Number($item.data('gzIndex'));
            return index;
        };

        this._create = function($comboItem){
            var comboWrapperWidth = $comboItem.outerWidth();
            this.element = $comboItem.attr('placeholder', defaults.placeholder).css({
                "width" : defaults.width,
                "line-height" : $comboItem.outerHeight()
            }).wrap('<div class="combo-control"></div>').after('<em class="combo-icon-arrow"></em>')
                .parents('div.combo-control').css({
                    'width': comboWrapperWidth + 'px',
                    "display" : defaults.displayMode,
                    "position" : "relative"
                });

            this.inputEl = $comboItem;

            if(!defaults.isAsync){
                this._generateDropList();
            }else{
                var $dropListWrapper = $('<ul class="gzCombo-drop-list"></ul>');

                this.dropList = $dropListWrapper.append('<div class="loader" style="display: none;"></div>')
                    .appendTo(this.element).hide();

                this.loader = this.dropList.find('>div.loader');
            }
            this._bindEvent();
        };

        this._generateDropList = function(){
            var $dropListWrapper = $('<ul class="gzCombo-drop-list"></ul>');

            this.dropList = $dropListWrapper.append('<div class="loader" style="display: none;"></div>')
                .appendTo(this.element).hide();

            this.loader = this.dropList.find('>div.loader');

            if(defaults.dataSource && defaults.dataSource instanceof Array && (defaults.dataSource.length > 0)){ //Array Data

                _processRenderData.call(this, defaults.dataSource);

                _processRenderedStyle.call(this);

            }else if(defaults.dataSource && typeof defaults.dataSource === 'string'){ //URL
                var deferredObj = $.ajax({
                    type: "POST",
                    url: defaults.dataSource,
                    data: defaults.requestData,
                    beforeSend: $.proxy(function(xhr){
                        this.loader.show();
                    }, this)
                });

                deferredObj.done($.proxy(function(data){
                    _processRenderData.call(this, data);

                    _processRenderedStyle.call(this);

                    this.loader.hide();

                }, this));
            }
        };

        this._bindEvent = function(){
            var $dropList = this.dropList,
                $document = $(document),
                keyMap = defaults.keyMap;

            this.inputEl.off('click').on('click', function(e){
                var evt = e || window.event;

                if($dropList.is(':hidden')){
                    if(defaults.isAsync && defaults.dataSource && typeof defaults.dataSource === 'string'){
                        if(!isLoaded){
                            var deferredObj = $.ajax({
                                type: "POST",
                                url: defaults.dataSource,
                                data: defaults.requestData,
                                beforeSend: function(xhr){
                                    gzComboPlugin.loader.show();
                                }
                            });

                            deferredObj.done(function(data){
                                _processRenderData.call(gzComboPlugin, data);

                                _processRenderedStyle.call(gzComboPlugin);

                                gzComboPlugin.loader.hide();
                                isLoaded = true;
                            });
                        }else{
                            $dropList.show();
                            gzComboPlugin._echoSelectedItem();
                        }
                    }else if(!defaults.isAsync){
                        $dropList.show();
                        gzComboPlugin._echoSelectedItem();
                    }
                }
                evt.stopPropagation();
                evt.preventDefault();
            }).on('keydown', function(e){
                var evt = e || window.event,
                    currentKey = evt.keyCode || evt.charCode || evt.which,
                    len = gzComboPlugin.dropList.children('.gzCombo-drop-list-item').length;

                if(gzComboPlugin.curDropListItem){
                    gzComboPlugin.curDropListItem.removeClass('gzCombo-over');
                }
               switch (currentKey){
                   case keyMap.ENTER:
                       gzComboPlugin._fillBack(gzComboPlugin.curDropListItem);
                       break;
                   case keyMap.UP:
                       if(keyIndex === -1){
                           gzComboPlugin.curDropListItem = gzComboPlugin.dropList.children('.gzCombo-drop-list-item').eq(len-1);
                           gzComboPlugin._addHoverClass(gzComboPlugin.curDropListItem);
                           keyIndex = len - 1;
                           return false;
                       }

                       if(keyIndex === 0){
                           gzComboPlugin.curDropListItem = gzComboPlugin.dropList.children('.gzCombo-drop-list-item').eq(len-1);
                           gzComboPlugin._addHoverClass(gzComboPlugin.curDropListItem);
                           keyIndex = len - 1;
                           return false;
                       }

                       gzComboPlugin.curDropListItem = gzComboPlugin.dropList.children('.gzCombo-drop-list-item').eq(--keyIndex);
                       gzComboPlugin._addHoverClass(gzComboPlugin.curDropListItem);
                       break;
                   case keyMap.DOWN:
                       if(keyIndex === len - 1){
                           gzComboPlugin.curDropListItem = gzComboPlugin.dropList.children('.gzCombo-drop-list-item').eq(0);
                           gzComboPlugin._addHoverClass(gzComboPlugin.curDropListItem);
                           keyIndex = 0;
                           return false;
                       }
                       gzComboPlugin.curDropListItem = gzComboPlugin.dropList.children('.gzCombo-drop-list-item').eq(++keyIndex);
                       gzComboPlugin._addHoverClass(gzComboPlugin.curDropListItem);
                       break;
               }
            });

            $document.off('click').on('click', function(e){
                $dropList.hide();
            });

            this.element.off('click', 'li.gzCombo-drop-list-item')
                .on('click', 'li.gzCombo-drop-list-item',function(e){
                var evt = e || window.event,
                    $self = $(this);

                gzComboPlugin._fillBack($self);
                evt.stopPropagation();
            }).on('mouseenter', 'li.gzCombo-drop-list-item',function(e){
                var $self = $(this);
                gzComboPlugin._addHoverClass($self);
            }).on('mouseleave', 'li.gzCombo-drop-list-item',function(e){
                var $self = $(this);
                gzComboPlugin._removeHoverClass($self);
            });
        };

        this._fillBack = function($dropListItem){
            var  labelValue = $dropListItem.data('labelValue'),
                labelName = $dropListItem.data('labelName');

            this.inputEl.val(labelName).data({
                'labelValue' : labelValue,
                'labelName' : labelName
            });

            this.dropList.hide();
        };

        this._echoSelectedItem = function(){
            this.dropList.children('.gzCombo-drop-list-item').each($.proxy(function(i, listItem){
                if($(listItem).data('labelValue') === this.inputEl.data('labelValue')){
                    $(listItem).addClass('gzCombo-current');
                    this.selectedItem = $(listItem);
                    keyIndex = this.getIndex($(listItem));
                }else{
                    $(listItem).removeClass('gzCombo-current');
                }
                this._removeHoverClass($(listItem));
            }, this));
        };

        this._addHoverClass = function($item){
            $item.addClass('gzCombo-over');
        };

        this._removeHoverClass = function($item){
            $item.removeClass('gzCombo-over');
        };

        return this.each($.proxy(function(i, comboItem){
            this._create($(comboItem));
        }, this));
    };

    $.fn.gzCombo = gzCombo;
})(window, jQuery);