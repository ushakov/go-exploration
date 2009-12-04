var tilelayer = null;
var myTileLayer = null;
var map = null;
var g_ob = new Object;
var edit = null;
var current_marker = null;
var icon_base = 'http://maps.google.com/mapfiles/ms/micons/';
var known_icons = ["red-dot.png", "blue-dot.png", "green-dot.png", "yellow-dot.png", "purple-dot.png"];

function initialize() {
    edit = document.getElementById("newtitle");
    // Set up the copyright information
    // Each image used should indicate its copyright permissions
    var myCopyright = new GCopyrightCollection("(c) ");
    myCopyright.addCopyright(new GCopyright('Demo',
                                            new GLatLngBounds(new GLatLng(-90,-180), new GLatLng(90,180)),
                                            0,'©2007 Google'));

    // Create the tile layer overlay and 
    // implement the three abstract methods   		
    tilelayer = new GTileLayer(myCopyright);
    tilelayer.getTileUrl = function(p, z) { return "/gettile?x=" + p.x +
                                               "&y=" + p.y + "&z=" + z; };
    tilelayer.isPng = function() { return true;};
    g_ob.opacity = 1.0;
    tilelayer.getOpacity = function() { return g_ob.opacity; }

    myTileLayer = new GTileLayerOverlay(tilelayer);
    map = new GMap2(document.getElementById("map_canvas"));
    map.setCenter(new GLatLng(55.75, 37.65), 13);
    map.addControl(new GSmallMapControl());
    map.addControl(new GMapTypeControl());
    map.enableScrollWheelZoom();
      
    map.addOverlay(myTileLayer);
}

function changeOpacity(delta) {
  opacity = g_ob.opacity + delta;
  if (opacity < 0.01) opacity = 0.01;
  if (opacity > 1.0) opacity = 1.0;
  g_ob.opacity = opacity;
  map.removeOverlay(myTileLayer);
  map.addOverlay(myTileLayer);
  return false;
}


