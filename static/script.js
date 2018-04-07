var canvas = document.getElementById('place_chain_canvas');
var canvascolor = document.getElementById('place_chain_color_chooser');
var colorindex = 0;
var pixelsize = 20;
var colors = [
	{r: 0, g: 0, b: 0},
	{r: 0, g: 0, b: 255},
	{r: 0, g: 255, b: 0},
	{r: 0, g: 255, b: 255},
	{r: 255, g: 0, b: 0},
	{r: 255, g: 0, b: 255},
	{r: 255, g: 255, b: 0},
	{r: 255, g: 255, b: 255},
	];
function getMousePos(canvas, evt) {
	var rect = canvas.getBoundingClientRect();
	return {
		x: evt.clientX - rect.left,
		y: evt.clientY - rect.top
	};
}
function roundPos(pos) {
	return {
		x: ((pos.x / pixelsize) | 0) * pixelsize,
		y: ((pos.y / pixelsize) | 0) * pixelsize,
	};
}
function setPixel(canvas, pos, color) {
	var context = canvas.getContext('2d');
	context.fillStyle = "rgb("+color.r+","+color.g+","+color.b+")";
	context.fillRect( pos.x, pos.y, pixelsize, pixelsize );
}
function refreshColorChooser(canvas) {
	var context = canvas.getContext('2d');
	for (j = 0; j < 2; j++) {
		for (i = 0; i < 4; i++) {
			color=colors[i+j*4];
			if(i+j*4 == colorindex){
				context.fillStyle = "grey"
				context.fillRect(i * 100,j * 100, 100, 100);
				context.fillStyle = "rgb(" + color.r + "," + color.g + "," + color.b + ")";
				context.fillRect(i * 100 + 5, j * 100 + 5, 90, 90);
			}else{
				context.fillStyle = "rgb(" + color.r + "," + color.g + "," + color.b + ")";
				context.fillRect(i * 100, j * 100, 100, 100);
			}
		}
	}
}
$("#place_chain_canvas").click(function(evt){
	var pos = roundPos(getMousePos(canvas, evt));
	pos.color = colorindex;
	setPixel(canvas,pos,colors[colorindex]);
	$.post("pixel",
		pos,
		function(data, status){
			$("#statusconsole").html("status = " + status + ", data = " + data);
		}
	);
});
$("#place_chain_color_chooser").click(function(evt){
		var mousePos = getMousePos(canvascolor, evt);
		colorindex = ((mousePos.x / 100) | 0) + ((mousePos.y / 100) | 0) * 4;
		refreshColorChooser(canvascolor);
});
refreshColorChooser(canvascolor);

//TODO longpoll node : image changes
//TODO longpoll node : my coins change
//TODO longpoll node : my color change ? oder nicht ?
