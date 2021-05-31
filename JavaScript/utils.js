/**
 * This JS-File contains some functions that are commonly used
 */

// The Utils "parent"
function Utils () {};

// Displays the given message of type "succes", "error" or "info"
Utils.showMessage = function(type, message) {
    switch (type) {
        case "success":
            alertify.success(message);
            break;
        case "error":
            alertify.error(message);
            break;
        case "info":
            alertify.warning(message);
    }
}

// Copies the given string to clipboard
Utils.copyToClipboard = function(text) {
    const el = document.createElement('textarea');
    el.value = text;
    el.setAttribute('readonly', '');
    el.style.position = 'absolute';
    el.style.left = '-9999px';
    document.body.appendChild(el);
    el.select();
    document.execCommand('copy');
    document.body.removeChild(el);
}

// Convert a single 0-F to 0-15
Utils.hex2num_single = function(hex) {
    if (hex < 10)
        return hex;
    switch(hex.toUpperCase()) {
        case "A":
            return 10;
        case "B":
            return 11;
        case "C":
            return 12;
        case "D":
            return 13;
        case "E":
            return 14;
        case "F":
            return 15;
    }
}

// Convert a single 0-15 to 0-F
Utils.num2hex_single = function(num) {
    if (num < 10)
        return num;
    switch(num) {
        case 10:
            return "A";
        case 11:
            return "B";
        case 12:
            return "C";
        case 13:
            return "D";
        case 14:
            return "E";
        case 15:
            return "F";
    }
}

// Converts a hex value to rgb
Utils.hexToRgb = function(hex) {
    var result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
    return result ? {
        r: parseInt(result[1], 16),
        g: parseInt(result[2], 16),
        b: parseInt(result[3], 16)
    } : null;
}

// Returns the browsers true width
Utils.getBrowserWidth = function() {
    return Math.max(
      document.body.scrollWidth,
      document.documentElement.scrollWidth,
      document.body.offsetWidth,
      document.documentElement.offsetWidth,
      document.documentElement.clientWidth
    );
}
  

// Removes any current drag selection (not supported on IE)
Utils.deleteSelection = function() {
    if (window.getSelection) {
        var selection = window.getSelection();
        selection.empty();
    }
}

// Scrolls to the destination in x miliseconds
Utils.scrollTo = function (final, duration) {
    var start = window.scrollY || document.documentElement.scrollTop,
        currentTime = null;
        
    var animateScroll = function(timestamp) {
        if (!currentTime) {
            currentTime = timestamp;  
        }      

        let progress = timestamp - currentTime;

        if(progress > duration) {
            progress = duration;
        }

        let val = Math.easeInOutQuad(progress, start, final-start, duration);
        window.scrollTo(0, val);

        if(progress < duration) {
            window.requestAnimationFrame(animateScroll);
        }
    };
  
    window.requestAnimationFrame(animateScroll);
};
  
// Deletes all stored cookies
Utils.deleteCookies = function() {
    var allCookies = document.cookie.split(';');
                
    for (var i = 0; i < allCookies.length; i++)
        document.cookie = allCookies[i] + "=;expires="
        + new Date(0).toUTCString();
}

// Used for animation curves
Math.easeInOutQuad = function (t, b, c, d) {
    t /= d/2;
    if (t < 1) return c/2*t*t + b;
    t--;
    return -c/2 * (t*(t-2) - 1) + b;
};