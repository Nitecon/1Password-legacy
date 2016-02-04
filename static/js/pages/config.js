/**
 * Created by whatting on 29/12/15.
 */
// Overall viewmodel for this screen, along with initial state
var dataUrl = "/rest/config";
var dm = {};

function PageViewModel(jsonData) {
    var self = this;
    self.ds = ko.viewmodel.fromModel(jsonData);
    self.validateVault = function(){
        ajaxAdd("/rest/config/validateVault", ko.toJSON(self.ds),function(){
                $.notifyBar({ cssClass: "success", html: "Path updated.  Please close the application and start it again." });
            $(".mainContainer").html('<div class="container"><div class="row"><h1>Application Updated, please restart!</h1></div></div>');
            }, function(d){
                $.notifyBar({ cssClass: "error", html: "Invalid location, path does not contain a '1password.html' file.<br>" });
            }
        );
    };
}

$(function() {
    ajaxGet(dataUrl, {},
        function(result) {
            dm = result;
            ko.applyBindings(new PageViewModel(result));
        },
        function() {
            $(".mainContainer").html("<h1>A problem occurred when trying to load the data.</h1>");
        });
});
