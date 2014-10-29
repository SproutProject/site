'use strict'

var qa = new function(){
    var that = this;

    that.load = function(){
	var template = $('#template').html();

	$.post('/spt/d/qa',{},function(res){
	    $('#list').html(Mustache.render(template,res));
	});
    };
}
