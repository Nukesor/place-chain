var canvas = document.getElementById('place_chain_canvas');
var canvascolor = document.getElementById('place_chain_color_chooser');
var pixelnumberdisplay = document.getElementById('pixelnumberdisplay');
var colorindex = 1;
var size = 10;
var lastdata;
var countPixels = 0;
var myprivkey;
var myname;
var mybio;
var myimg;
var loggedin = 0;
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

function toHexString(byteArray) {
  return Array.from(byteArray, function(byte) {
    return ('0' + (byte & 0xFF).toString(16)).slice(-2);
  }).join('')
}

function setPixel(canvas, pos, color) {
	var rect = canvas.getBoundingClientRect();
	var context = canvas.getContext('2d');
	context.fillStyle = "rgb(" + color.r + "," + color.g + "," + color.b + ")";
	context.fillRect(pos.x * rect.width / size , pos.y * rect.height / size, rect.width / size, rect.height / size);
}

function refreshColorChooser(canvas) {
	var context = canvas.getContext('2d');
	for (i = 0; i < 2; i++) {
		for (j = 0; j < 4; j++) {
			color = colors[i * 4 + j + 1];
			if(colorindex != 0) {
				if(i * 4 + j + 1 == colorindex) {
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
	$("#p_name_owner").text(lastdata[pos.x][pos.y].Profile.Name);
	$("#p_bio_owner").text(""+lastdata[pos.x][pos.y].Profile.Bio);
	$("#p_image_owner").attr("src", lastdata[pos.x][pos.y].Profile.AvatarUrl);
	$("#p_image_owner").attr("src", "https://upload.wikimedia.org/wikipedia/commons/4/4d/CowsGoMoOo_%28Alexander_Huard%29.png");
});

$("#place_chain_canvas").click(function(evt) {
	if(loggedin>0){
		var pos = getMousePos(canvas, evt);
		setPixel(canvas, pos, colors[colorindex]);
		var pixel = {X: pos.x, Y: pos.y, Color: colorindex, Nonce: Math.random().toString(36).substring(7)};
		var pixelstring = JSON.stringify(pixel);
		var signature = myprivkey.signString(pixelstring);
		pixel.PubKey = { type: "ed25519", data: toHexString(myprivkey.makePubKey().bytes)}
		pixel.Signature = { type: "ed25519", data: toHexString(signature.bytes)}
		$.ajax("pixel", {
			data : JSON.stringify(pixel),
			contentType : 'application/json',
			type : 'POST',})
		.done(function(msg) {$("#statusconsole").html("msg = done");})
		.fail(function(xhr, status, error) {$("#statusconsole").html("msg = " + status + ", error = " + error);});
	}
});

$("#place_chain_color_chooser").click(function(evt) {
	var rect = canvascolor.getBoundingClientRect();
	var mousePos = {
		x: ((evt.clientX - rect.left) / rect.width * 2) | 0,
		y: ((evt.clientY - rect.top) / rect.height * 4) | 0
	};
	colorindex = mousePos.x * 4 + mousePos.y + 1;
	refreshColorChooser(canvascolor);
});

$("#register_button").click(function(evt) {
	myname = $("#name_input").val();
	mybio = $("#bio_input").val();
	myimg = $("#img_input").val();
	if(myimg == ""){
		myimg = "blank_profile_100.png"
	}
	if(myname == ""){
		myname = Math.random().toString(36).substring(7)
	}
	if(mybio == ""){
		mybio = Math.random().toString(36).substring(7)
	}
	myprivkey = tendermintcrypto.genPrivKeyEd25519();
	var data = {
		PubKey: { type: "ed25519", data: toHexString(myprivkey.makePubKey().bytes)},
		Profile: {Name: myname, Bio: mybio, AvatarUrl: myimg}
	};
	$.ajax("register", {
		data : JSON.stringify(data),
		contentType : 'application/json',
		type : 'POST'})
	.done(function(msg) {
		$("#p_name").text(myname);
		$("#p_bio").text(mybio);
		$("#p_privkey").text(toHexString(myprivkey.bytes));
		$("#p_image").attr("src", myimg);
		$("#loginregister_div").hide();
		$("#profile_div").show();
		loggedin = 1;
	})
	.fail(function(xhr, status, error) {
		$("#statusconsole").html("msg = " + status + ", error = " + error);
	});
});

$("#login_button").click(function(evt) {
	myname = $("#name_input").val();
	mybio = $("#bio_input").val();
	myimg = $("#img_input").val();
	if(myimg == ""){
		myimg = "blank_profile_100.png"
	}
	myprivkey = tendermintcrypto.genPrivKeyEd25519();
	$("#p_name").text(myname);
	$("#p_bio").text(mybio);
	$("#p_privkey").text(toHexString(myprivkey.bytes));
	$("#p_image").attr("src", myimg);
	$("#loginregister_div").hide();
	$("#profile_div").show();
	loggedin = 1;
});

$(function() {
	$("#profile_div").hide();
	refreshColorChooser(canvascolor);
	setInterval(function() {
		$.get("pixels", function(data) {
			lastdata = data;
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
					if(data[i][j].Color != 0) {
						setPixel(canvas,{x: i, y: j}, colors[data[i][j].Color]);
						countPixels = countPixels +1;
					}
				}
			}
			$("#pixelnumberdisplay").html("Number of Pixels: " + countPixels);
			countPixels = 0;
		});
	}, 1000);
	//TODO longpoll node : image changes
	//TODO longpoll node : my coins change
	//TODO longpoll node : my color change ? oder nicht ?
});
