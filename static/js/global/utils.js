/**
 * Created by whatting on 9/16/15.
 */

$(function () {
    $.ajaxSetup({ cache: false });
});

function ajaxGet(url,dataTosave,callback,errCallback) {
    ajaxModify(url, dataTosave, "GET", callback, errCallback);
}

function ajaxAdd(url, dataTosave, callback, errCallback) {
    ajaxModify(url, dataTosave, "POST", callback, errCallback);
}

function ajaxUpdate(url, dataToSave, successCallback, errCallback){
    ajaxModify(url, dataToSave, "PUT", successCallback, errCallback);
}

function ajaxDelete(url, callback, errCallback) {
    ajaxModify(url, null, "DELETE", callback, errCallback);
}

function ajaxModify(url, dataToSave, httpMethod, callback, errCallback){
    $.ajax({
            type: httpMethod,
            url: url,
            data:dataToSave,
            contentType:"application/json",
            dataType: "json"
        })
        .done(function(data){
            /* Add console.log stuff here if need be */
            if (callback !== undefined){
                callback(data);
            }
        })
        .fail(function(data){
            if (errCallback !== undefined){
                errCallback(data);
            }
        });
}

var SlideItemDown = function(d,ev) {
    $(ev.target).hide().slideDown();
};

var SlideItemUpAndRemove = function(d,ev) {
    $(ev.target).slideUp(function() { $(element).remove(); });
};

var DoPopOver = function(data,ev){
    $(function() {
        console.log($(ev.target).advancedpopover());
        console.log(ev.target);
    });
};


ko.bindingHandlers.executeOnEnter = {
    init: function (element, valueAccessor, allBindingsAccessor, viewModel){
        var data = valueAccessor();
        $(element).keypress(function(event){
            var keyCode = (event.which ? event.which : event.keyCode);
            if (keyCode ===13){
                data.call(viewModel);
                return false;
            }
            return true;
        });
    }
};

ko.bindingHandlers.executeOnClick = {
    init: function (element, valueAccessor, allBindingsAccessor, viewModel){
        var data = valueAccessor();
        $(element).click(function(event){
            data.call(viewModel);
            return true;
        });
    }
};

ko.bindingHandlers.executeOnBlur = {
    init: function (element, valueAccessor, allBindingsAccessor, viewModel){
        var data = valueAccessor();
        $(element).blur(function(event){
            data.call(viewModel);
            return true;
        });
    }
};

ko.bindingHandlers.modal = {
    init: function (element, valueAccessor, allBindingsAccessor, viewModel, bindingContext) {
        $(element).keypress(function(event){
            var keyCode = (event.which ? event.which : event.keyCode);
            if (keyCode ===13){
                valueAccessor().call(viewModel);
                return false;
            }
            return true;
        });
        $(element).modal({ show: false }).on("hidden", function () {
            var data = valueAccessor();
            if (ko.isWriteableObservable(data))
                data(null);
        });



        ko.applyBindingsToNode(element, { with: valueAccessor() }, bindingContext);

        return { controlsDescendantBindings: true };
    },
    update: function (element, valueAccessor) {
        var data = ko.unwrap(valueAccessor());

        $(element).modal( data ? "show" : "hide" );
    }
};

ko.bindingHandlers.beforeUnload = {
    init: function(element, valueAccessor, allBindingsAccessor, viewModel) {
        if (window.onbeforeunload == null) {
            window.onbeforeunload = function(){
                var value = valueAccessor();
                value.call(viewModel);
                return undefined;
            };

        } else {
            var err = "onbeforeupload has already been set";
            throw new Error(err);
        }
    }
};

/* catch and precent cyclical references in rootObject */
function toJSON(rootObject, replacer, spacer) {
    var cache = [];
    var plainJavaScriptObject = ko.toJS(rootObject);
    var replacerFunction = replacer || cycleReplacer;
    var output = ko.utils.stringifyJson(plainJavaScriptObject, replacerFunction, spacer || 2);
    cache = null;
    return output;

    function cycleReplacer(key, value) {
        if (typeof value === 'object' && value !== null) {
            if (cache.indexOf(value) !== -1) {
                return; // cycle is found, skip it
            }
            cache.push(value);
        }
        return value;
    }
}

ko.bindingHandlers.dump = {
    init: function(element, valueAccessor, allBindingsAccessor, viewmodel, bindingContext) {
        var context = valueAccessor();
        var allBindings = allBindingsAccessor();
        var pre = document.createElement('pre');

        element.appendChild(pre);

        var dumpJSON = ko.computed({
            read: function() {
                var enable = allBindings.enable === undefined || allBindings.enable;
                return enable ? toJSON(context, null, 2): '';
            },
            disposeWhenNodeisRemoved: element
        });

        ko.applyBindingsToNode(pre, {
            text: dumpJSON,
            visible: dumpJSON
        });

        return { controlsDescendentBindings: true };
    }
};

function formatDate(inputFormat) {
    var date = new Date(inputFormat);
    return [(date.getMonth() + 1), (date.getDate()), date.getFullYear()].join('/');
}

/*
 * Notify Bar - jQuery plugin
 *
 * Copyright (c) 2009-2015 Dmitri Smirnov
 *
 * Licensed under the MIT license:
 * http://www.opensource.org/licenses/mit-license.php
 *
 * Project home:
 * http://www.whoop.ee/posts/2013/04/05/the-resurrection-of-jquery-notify-bar.html
 * https://github.com/dknight/jQuery-Notify-bar
 *
 * Example Usage:
 *   $.notifyBar();
 *   $.notifyBar({ cssClass: "error", html: "Error occurred!" });
 *   $.notifyBar({ cssClass: "success", html: "Your data has been changed!" });
 *   $.notifyBar({ cssClass: "warning", html: "Settings aren't changed!" });
 *   $.notifyBar({ cssClass: "custom", html: "Your data has been changed!" });
 *   $.notifyBar({ html: "Click 'close' to hide notify bar", close: true, closeOnClick: false });
 *   $.notifyBar({ html: "At bottom", position: "bottom" });
 */
(function($) {

    "use strict";

    $.notifyBar = function(options) {
        var rand = parseInt(Math.random() * 100000000, 0),
            text_wrapper, asTime,
            bar = {},
            settings = {};

        settings = $.extend({
            html: 'Your message here',
            delay: 3000,
            animationSpeed: 200,
            cssClass: '',
            jqObject: '',
            close: false,
            closeText: '&times;',
            closeOnClick: true,
            closeOnOver: false,
            onBeforeShow: null,
            onShow: null,
            onBeforeHide: null,
            onHide: null,
            position: 'top'
        }, options);

        // Use these methods as private.
        this.fn.showNB = function() {
            if (typeof settings.onBeforeShow === 'function') {
                settings.onBeforeShow.call();
            }
            $(this).stop().slideDown(asTime, function() {
                if (typeof settings.onShow === 'function') {
                    settings.onShow.call();
                }
            });
        };

        this.fn.hideNB = function() {
            if (typeof settings.onBeforeHide === 'function') {
                settings.onBeforeHide.call();
            }
            $(this).stop().slideUp(asTime, function() {
                if (bar.attr("id") === "__notifyBar" + rand) {
                    $(this).slideUp(asTime, function() {
                        $(this).remove();
                        if (typeof settings.onHide === 'function') {
                            settings.onHide.call();
                        }
                    });
                } else {
                    $(this).slideUp(asTime, function() {
                        if (typeof settings.onHide === 'function') {
                            settings.onHide.call();
                        }
                    });
                }
            });
        };

        if (settings.jqObject) {
            bar = settings.jqObject;
            settings.html = bar.html();
        } else {
            bar = $("<div></div>")
                .addClass("jquery-notify-bar")
                .addClass(settings.cssClass)
                .attr("id", "__notifyBar" + rand);
        }
        text_wrapper = $("<span></span>")
            .addClass("notify-bar-text-wrapper")
            .html(settings.html);

        bar.html(text_wrapper).hide();

        var id = bar.attr("id");
        switch (settings.animationSpeed) {
            case "slow":
                asTime = 600;
                break;
            case "default":
            case "normal":
                asTime = 400;
                break;
            case "fast":
                asTime = 200;
                break;
            default:
                asTime = settings.animationSpeed;
        }
        $("body").prepend(bar);

        // Style close button in CSS file
        if (settings.close) {
            // If close settings is true. Set delay to one billion seconds.
            // It'a about 31 years - mre than enough for cases when notify bar is used :-)
            settings.delay = Math.pow(10, 9);
            bar.append($("<a href='#' class='notify-bar-close'>" + settings.closeText + "</a>"));
            $(".notify-bar-close").click(function(event) {
                event.preventDefault();
                bar.hideNB();
            });
        }

        // Check if we've got any visible bars and if we have,
        // slide them up before showing the new one
        if ($('.jquery-notify-bar:visible').length > 0) {
            $('.jquery-notify-bar:visible').stop().slideUp(asTime, function() {
                bar.showNB();
            });
        } else {
            bar.showNB();
        }

        // Allow the user to click on the bar to close it
        if (settings.closeOnClick) {
            bar.click(function() {
                bar.hideNB();
            });
        }

        // Allow the user to move mouse on the bar to close it
        if (settings.closeOnOver) {
            bar.mouseover(function() {
                bar.hideNB();
            });
        }

        setTimeout(function() {
            bar.hideNB(settings.delay);
        }, settings.delay + asTime);

        if (settings.position === 'bottom') {
            bar.addClass('bottom');
        } else if (settings.position === 'top') {
            bar.addClass('top');
        }

        return bar;
    };
})(jQuery);