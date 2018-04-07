var canvas = document.getElementById('place_chain_canvas');
var canvascolor = document.getElementById('place_chain_color_chooser');
var colorindex = 1;
var size = 10;
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
		x: ((evt.clientX - rect.left) / rect.width * size) | 0,
		y: ((evt.clientY - rect.top) / rect.height * size) | 0
	};
}
function setPixel(canvas, pos, color) {
	var rect = canvas.getBoundingClientRect();
	var context = canvas.getContext('2d');
	context.fillStyle = "rgb(" + color.r + "," + color.g + "," + color.b + ")";
	context.fillRect(pos.x * rect.width / size , pos.y * rect.height / size, rect.width / size, rect.height / size);
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
	var pos = getMousePos(canvas, evt);
	var pixel = {x: pos.x, y: pos.y, color: colorindex};
	$("#coordinates").html("x=" + pixel.x + ", y=" + pixel.y);
});
$("#place_chain_canvas").click(function(evt) {
	var pos = getMousePos(canvas, evt);
	setPixel(canvas, pos, colors[colorindex]);
	var pixel = {x: pos.x, y: pos.y, color: colorindex};
	$.ajax("pixel", {
        data : JSON.stringify(pixel),
        contentType : 'application/json',
        type : 'POST',
    }).done(function(msg) {$("#statusconsole").html("msg = " + msg);})
      .fail(function(xhr, status, error) {$("#statusconsole").html("msg = " + status + ", error = " + error);});
});
$("#place_chain_color_chooser").click(function(evt) {
	var rect = canvascolor.getBoundingClientRect();
	var mousePos = {
		x: ((evt.clientX - rect.left) / rect.width * 4) | 0,
		y: ((evt.clientY - rect.top) / rect.height * 2) | 0
	};
	colorindex = mousePos.x + mousePos.y * 4 + 1;
	refreshColorChooser(canvascolor);
});
$(function() {
	refreshColorChooser(canvascolor);
	$.get("pixels", function(data) {
		data = JSON.parse(data);
		width = data.length;
console.log("a");
		if(width == 0) {
			$("#statusconsole").html("error: requesting \"pixels\" returned zero width");
			return;
		}
console.log("b");
console.log(width);
		height = data[0].length;
console.log(height);
console.log(data);
		if(width != height) {
			$("#statusconsole").html("error: requesting \"pixels\" returned not square area");
			return;
		}
		size = width;
console.log(size);
		for(i = 0; i < size; i++) {
			for(j = 0; j < size; j++) {
				if(colors[data[i][j]] != 0) {
					setPixel(canvas,{x: i, y: j}, colors[data[i][j]]);
				}
			}
		}
	});
	//TODO longpoll node : image changes
	//TODO longpoll node : my coins change
	//TODO longpoll node : my color change ? oder nicht ?
});

