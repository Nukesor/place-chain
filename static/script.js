var canvas = document.getElementById('place_chain_canvas');
var canvascolor = document.getElementById('place_chain_color_chooser');
var colorindex = 1;
var pixelsize = 20;
var colors = [
	{},			//error color
	{r: 0, g: 0, b: 0}, 	//black
	{r: 0, g: 0, b: 255}, 	//red
	{r: 0, g: 255, b: 0},	//blue
	{r: 0, g: 255, b: 255},	// purple
	{r: 255, g: 0, b: 0}, 	// green
	{r: 255, g: 0, b: 255},	// yellow
	{r: 255, g: 255, b: 0},	// bright blue
	{r: 255, g: 255, b: 255},	// white
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
	context.fillStyle = "rgb(" + color.r + "," + color.g + "," + color.b + ")";
	context.fillRect( pos.x, pos.y, pixelsize, pixelsize );
}
function refreshColorChooser(canvas) {
	var context = canvas.getContext('2d');
	for (j = 0; j < 2; j++) {
		for (i = 0; i < 4; i++) {
			color = colors[i + j * 4 + 1];
			if(colorindex != 0) {
				if(i + j * 4 + 1 == colorindex) {
					context.fillStyle = "grey"
					context.fillRect(i * 100, j * 100, 100, 100);
					context.fillStyle = "rgb(" + color.r + "," + color.g + "," + color.b + ")";
					context.fillRect(i * 100 + 5, j * 100 + 5, 90, 90);
				} else {
					context.fillStyle = "rgb(" + color.r + "," + color.g + "," + color.b + ")";
					context.fillRect(i * 100, j * 100, 100, 100);
				}
			}
		}
	}
}
$("#place_chain_canvas").mousemove(function(evt) {
	var pos = roundPos(getMousePos(canvas, evt));
	var pixel = {x: pos.x / pixelsize, y: pos.y / pixelsize, color: colorindex};
	$("#coordinates").html("x=" + pixel.x + ", y=" + pixel.y);
});
$("#place_chain_canvas").click(function(evt) {
	var pos = roundPos(getMousePos(canvas, evt));
	setPixel(canvas,pos,colors[colorindex]);
	var pixel = {x: pos.x / pixelsize, y: pos.y / pixelsize, color: colorindex};
	$.post("pixel", pixel)
		.done(function(msg) {$("#statusconsole").html("msg = " + msg);})
		.fail(function(xhr, status, error) {$("#statusconsole").html("msg = " + status + ", error = " + error);});
});
$("#place_chain_color_chooser").click(function(evt) {
		var mousePos = getMousePos(canvascolor, evt);
		colorindex = ((mousePos.x / 100) | 0) + ((mousePos.y / 100) | 0) * 4 + 1;
		refreshColorChooser(canvascolor);
});
$(function() {
	refreshColorChooser(canvascolor);
	$.get("pixels", function(data) {
		width = data.length;
		if(width == 0) {
			$("#statusconsole").html("error: requesting \"pixels\" returned zero width");
			return;
		}
		height = data[0].length;
		if(width != height) {
			$("#statusconsole").html("error: requesting \"pixels\" returned not square area");
			return;
		}
		size = width;
		for(i = 0; i < size; i++) {
			for(j = 0; j < size; j++) {
				setPixel(canvas,{x: i * pixelsize, y: j * pixelsize}, colors[data[i][j]]);
			}
		}
	});
	//TODO longpoll node : image changes
	//TODO longpoll node : my coins change
	//TODO longpoll node : my color change ? oder nicht ?
});

