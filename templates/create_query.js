var data = {};

function submit_search() {
    var genre_list  = document.getElementById("genre").children;
    var rating_list = document.getElementById("rating").children;
    var imdb_list   = document.getElementById("imdb").children;

    data["srch"]   = document.getElementById("search-query").value;
    data["genre"]  = get_selected(genre_list);
    data["rating"] = get_selected(rating_list);
    data["imdb"]   = get_selected(imdb_list);
    console.log(data);

    var URL = construct_url(data);
    console.log(URL);
    return URL;
    /* So this was a stupid idea but worth remembering

    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/results", true);
    xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
    xhr.send(JSON.stringify(data));

    xhr.onloadend = function () {
	// stuff to do after it sends?
    };
    */
    /* Backup Route (SIN)

    var new_form = document.createElement("form");
    new_form.method = "POST";
    new_form.action = "/results";
    new_form.style.visibility = "hidden";
    document.getElementById("main").appendChild(new_form);

    var query_string = JSON.stringify(data);
    var new_input = document.createElement("input");
    new_input.type = "text";
    new_input.name = "query_string";
    new_input.value = query_string;
    new_input.style.visibility = "hidden";
    new_form.appendChild(new_input);
    new_form.submit();
    */
}

// Creates the url
function construct_url(data) {
    var URL     = "/results?";
    var srch_q  = data['srch'];
    var genres  = data['genre'];
    var ratings = data['rating'];
    var imdb    = data['imdb'];

    console.log(srch_q + "\n ^ that should be the search box");
    URL += "q=" + encodeURIComponent(srch_q) + "&";
    

    var srch_pfix;
    srch_pfix = "gen="; // Genres
    URL = construct_url_helper(URL, genres, srch_pfix);
    URL += "&";
    console.log(genres);
    console.log(URL);

    srch_pfix = "rating="; // MPAA Ratings
    URL = construct_url_helper(URL, ratings, srch_pfix);
    URL += "&";
    console.log(ratings);
    console.log(URL);

    srch_pfix = "imdb="; // IMDb Ratings
    URL = construct_url_helper(URL, imdb, srch_pfix);
    console.log(imdb);
    console.log(URL);

    return URL;
}
// helper function
function construct_url_helper(URL, arr, pfix) {
    var res = URL; // stores result
    if (arr.indexOf("All") > -1) {
	// Just search "All" films... This case shouldn't happen, but this is a safeguard.
	return res + pfix;
    } else {
	for (var i = 0; i < arr.length; i++) {
	    // 
	    res += pfix + arr[i];
	    if (i !== arr.length - 1) {
		res += "&";
	    }
	}
    }
    return res;
}
// Removes .highlight class from the "All" button of a given unordered list (node)
function deselect_All(ul) {
    var All_btn = ul.children[0].children[0];
    if (All_btn.className.indexOf(" highlight") > -1) {
	All_btn.className = All_btn.className.replace(" highlight","");
    }
}

// Add .highlight class to the "All" button of a given unordered list (node)
function select_All(ul) {
    var All_btn = ul.children[0].children[0];
    if (All_btn.className.indexOf(" highlight") == -1) {
	All_btn.className += " highlight";
    }
}

// Removes .highlight class from any non-"All" buttons 
function deselect_others(ul) {
    var lis = ul.children;
    for (var i = 0; i < lis.length; i++) {
	var btn = lis[i].children[0];
	if (btn) {
	    if (btn.innerText !== "All") {
		btn.className = btn.className.replace(" highlight", "");
	    }
	}
    }
}

// adds/removes .highlight class from a button
// has some conditions to keep highlighting consistent
// this is a total mess, I'm sorry if you're reading this after the fact
// triggered onclick
function select_field(butt) {
    var b = butt.className;
    if (b.indexOf(" highlight") > -1) {
	butt.className = b.replace(" highlight","");
	var selected_btns = get_selected(butt.parentNode.parentNode.children);
	console.log(selected_btns);
	if (selected_btns.length === 1 && selected_btns[0] === "All") {
	    console.log("met condition");
	    select_All(butt.parentNode.parentNode);
	}
    } else {
	butt.className += " highlight";
	if (butt.innerText !== "All") {
	    deselect_All(butt.parentNode.parentNode);
	} else {
	    deselect_others(butt.parentNode.parentNode);
	}
    }
}

// get all selected buttons (buttons with .highlight class)
// returns Array
function get_selected(node_list) {
    var res = [];   // stores result
    for (var i = 0; i < node_list.length; i++) {
	var li = node_list[i];
	var btn = li.children;
	if (btn.length > 0) {
	    if (btn[0].className.indexOf("highlight") > -1) {
		res.push(btn[0].innerText)
	    }
	}
    }
    if (res.length === 0) {
	res.push("All")
    }
    return res;
}