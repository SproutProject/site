'use strict'

var poll = new function(){
    var that = this;

    that.load = function(){
	var template = $('#template').html();

	$.post('/spt/d/poll',{},function(res){
	    $('#list').html(Mustache.render(template,res));

	    $('#list div.item').on('click',function(e){
		var j_this = $(this);
		$('#list').find('div.subject').hide();
		j_this.addClass('active');
	    });
	    $('#list button.close').on('click',function(e){
		var j_item = $(this).parents('div.item');
		$('#list').find('div.subject').show();
		j_item.removeClass('active');
		return false;
	    });
	});
    };
}
