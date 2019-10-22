// Parallel Coordinates
// Released under the BSD License: http://opensource.org/licenses/BSD-3-Clause

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// INIT
///////////////////////////////////////////////////////////////////////////////////////////////////////////////

let width = document.body.clientWidth,
    height = d3.max([document.body.clientHeight-540, 240]);

let m = [60, 0, 10, 0],
    w = width - m[1] - m[3],
    h = height - m[0] - m[2],
    xscale = d3.scale.ordinal().rangePoints([0, w], 1),
    yscale = {},
    dragging = {},
    line = d3.svg.line(),
    axis = d3.svg.axis().orient("left").ticks(1+height/50),
    data,
    foreground,
    background,
    highlighted,
    dimensions,
    legend,
    render_speed = 50,
    brush_count = 0,
    excluded_groups = [];

// hsla
let colors = {
    "cDAI": [39,99,63],
    "cUSDC": [206,100,41],
    "cETH": [155,100,37],
    "cBAT": [14,100,50],
    "cZRX": [53,100,50],
    "cREP": [270,49,40],
    "cWBTC": [34,100,50],
};

// Scale chart and canvas height
d3.select("#chart")
    .style("height", (h + m[0] + m[2]) + "px")

d3.selectAll("canvas")
    .attr("width", w)
    .attr("height", h)
    .style("padding", m.join("px ") + "px");

// Foreground canvas for primary view
foreground = document.getElementById('foreground').getContext('2d');
foreground.globalCompositeOperation = "destination-over";
foreground.strokeStyle = "rgba(0,100,160,0.1)";
foreground.lineWidth = 1.7;
foreground.fillText("Loading...",w/2,h/2);

// Highlight canvas for temporary interactions
highlighted = document.getElementById('highlight').getContext('2d');
highlighted.strokeStyle = "rgba(0,100,160,1)";
highlighted.lineWidth = 4;

// Background canvas
background = document.getElementById('background').getContext('2d');
background.strokeStyle = "rgba(0,100,160,0.1)";``
background.lineWidth = 1.7;

// SVG for ticks, labels, and interactions
let svg = d3.select("svg")
    .attr("width", w + m[1] + m[3])
    .attr("height", h + m[0] + m[2])
    .append("svg:g")
    .attr("transform", "translate(" + m[3] + "," + m[0] + ")");


///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// LOAD DATA & VISUALIZATION
///////////////////////////////////////////////////////////////////////////////////////////////////////////////

let sources = {
    all: "all",
    all2: "all_2",
    all3: "all_3",
    cBAT: "cBAT",
    cDAI: "cDAI",
    cETH: "cETH",
    cREP: "cREP",
    cWBTC: "cWBTC",
    cZRX: "cWBTC",
}
let source = sources.all3;

d3.csv(`compound_${source}.csv`, function(raw_data) {

    console.log("raw_data:", raw_data);

    // Convert quantitative scales to floats
    data = raw_data.map(function(d) {
        // console.log("d:", d);
        for (let k in d) {
            if (!_.isNaN(raw_data[0][k] - 0) && k != "address") {
                d[k] = parseFloat(d[k]) || 0;
            }
        };
        return d;
    });

    // data = data.filter(d => {
    //   // console.log("d", d);
    //   return d.cToken === "cDai"});

    //   console.log("data", data);

    // Extract the list of numerical dimensions and create a scale for each.
    xscale.domain(dimensions = d3.keys(data[0]).filter(function(k) {
        // console.log("k:", k);
        switch(k) {
            case "address":
                return (_.isNumber(data[0][k])) && (yscale[k] = d3.scale.ordinal().rangePoints([0, height]));
            case "type":
                return (_.isNumber(data[0][k])) && (yscale[k] = d3.scale.ordinal().rangePoints([0, height]));
            case "cToken":
                return (_.isNumber(data[0][k])) && (yscale[k] = d3.scale.ordinal().rangePoints([0, height]));
            default:
                return (_.isNumber(data[0][k])) && (yscale[k] = d3.scale.sqrt()
                    .domain(d3.extent(data, function(d) { return +d[k]; }))
                    .range([h, 0]));
        }
    }));

    console.log("xscale:", xscale);
    console.log("dimensions before:", dimensions);


    /*

    Descriptive Stats
    Money Markets
    - tallies for each metric
       - sum
       - min/max
       - mean
       - std deviations / quartiles
       - distributions



    TODO :
    Filter 0 values toggle per dimension
    accurate update calls for tally
    filter brush works on min and max
    filter brush (min/max) value input


    */
    dimensionsObjects = [
        // {
        //   name: "address",
        //   scale: d3.scale.ordinal().rangePoints([0, height]),
        //   type: "string"
        // },
        // {
        //   name: "type",
        //   scale: d3.scale.ordinal().rangePoints([0, height]),
        //   type: "string"
        // },
        // {
        //   name: "cToken",
        //   scale: d3.scale.ordinal().rangePoints([0, height]),
        //   type: "string"
        // },

        {
            name: "total_minted",
            scale: d3.scale.linear().range([height, 0]),
            type: "number"
        },
        {
            name: "supply_balance",
            scale: d3.scale.linear().range([height, 0]),
            type: "number"
        },
        {
            name: "total_redeemed",
            scale: d3.scale.linear().range([height, 0]),
            type: "number"
        },
        {
            name: "number_of_mints",
            scale: d3.scale.linear().range([height, 0]),
            type: "number"
        },
        {
            name: "number_of_redeems",
            scale: d3.scale.linear().range([height, 0]),
            type: "number"
        },


        {
            name: "total_borrowed",
            scale: d3.scale.linear().range([height, 0]),
            type: "number"
        },
        {
            name: "total_repaid",
            scale: d3.scale.linear().range([height, 0]),
            type: "number"
        },
        {
            name: "borrow_balance",
            scale: d3.scale.linear().range([height, 0]),
            type: "number"
        },
        {
            name: "number_of_borrows",
            scale: d3.scale.linear().range([height, 0]),
            type: "number"
        },
        {
            name: "number_of_repayments",
            scale: d3.scale.linear().range([height, 0]),
            type: "number"
        },

        {
            name: "number_of_repayments",
            scale: d3.scale.linear().range([height, 0]),
            type: "number"
        },
        // liquidators
        {
            name: "total_liquidated",
            scale: d3.scale.linear().range([height, 0]),
            type: "number"
        },
        {
            name: "number_of_liquidations",
            scale: d3.scale.linear().range([height, 0]),
            type: "number"
        },
        // liquidatees
        {
            name: "number_of_liquidated",
            scale: d3.scale.linear().range([height, 0]),
            type: "number"
        },
        {
            name: "total_liquidated_for",
            scale: d3.scale.linear().range([height, 0]),
            type: "number"
        },
        // cToken Transfers
        {
            name: "number_of_transfers_in",
            scale: d3.scale.linear().range([height, 0]),
            type: "number"
        },
        {
            name: "total_transfers_in",
            scale: d3.scale.linear().range([height, 0]),
            type: "number"
        },
        {
            name: "number_of_transfers_out",
            scale: d3.scale.linear().range([height, 0]),
            type: "number"
        },
        {
            name: "total_transfers_out",
            scale: d3.scale.linear().range([height, 0]),
            type: "number"
        },
        {
            name: "num_actions",
            scale: d3.scale.linear().range([height, 0]),
            type: "number"
        },
        // dates
        // first_action,
        // last_action,
        {
            name: "days_active",
            scale: d3.scale.linear().range([height, 0]),
            type: "number"
        },
        {
            name: "days_since_first",
            scale: d3.scale.linear().range([height, 0]),
            type: "number"
        },
        {
            name: "days_since_last",
            scale: d3.scale.linear().range([height, 0]),
            type: "number"
        },
        // other
        {
            name: "ETH_borrowing_power",
            scale: d3.scale.linear().range([height, 0]),
            type: "number"
        },
        {
            name: "supply_balance_ETH",
            scale: d3.scale.linear().range([height, 0]),
            type: "number"
        },
        {
            name: "borrow_balance_ETH",
            scale: d3.scale.linear().range([height, 0]),
            type: "number"
        },
        {
            name: "collateral_ratio",
            scale: d3.scale.linear().range([height, 0]),
            type: "number"
        },
        // based on
        {
            name: "percent_redeemable",
            scale: d3.scale.linear().range([height, 0]),
            type: "number"
        }
    ];

    dimensions = dimensionsObjects.map(d => d.name);

    console.log("dimensions after:", dimensions);

    dimensionsObjects.forEach(function(dimensionObject) {
        dimensionObject.scale.domain(dimensionObject.type === "number"
            ? d3.extent(data, function(d) { return +d[dimensionObject.name]; })
            : data.map(function(d) { return d[dimensionObject.name]; }).sort());
    });

    // Add a group element for each dimension.
    let g = svg.selectAll(".dimension")
        .data(dimensions)
        .enter().append("svg:g")
        .attr("class", "dimension")
        .attr("transform", function(d) { return "translate(" + xscale(d) + ")"; })
        .call(d3.behavior.drag()
            .on("dragstart", function(d) {
                dragging[d] = this.__origin__ = xscale(d);
                this.__dragged__ = false;
                d3.select("#foreground").style("opacity", "0.35");
            })
            .on("drag", function(d) {
                dragging[d] = Math.min(w, Math.max(0, this.__origin__ += d3.event.dx));
                dimensions.sort(function(a, b) { return position(a) - position(b); });
                xscale.domain(dimensions);
                g.attr("transform", function(d) { return "translate(" + position(d) + ")"; });
                brush_count++;
                this.__dragged__ = true;

                // Feedback for axis deletion if dropped
                if (dragging[d] < 12 || dragging[d] > w-12) {
                    d3.select(this).select(".background").style("fill", "#b00");
                } else {
                    d3.select(this).select(".background").style("fill", null);
                }
            })
            .on("dragend", function(d) {
                if (!this.__dragged__) {
                    // no movement, invert axis
                    let extent = invert_axis(d);

                } else {
                    // reorder axes
                    d3.select(this).transition().attr("transform", "translate(" + xscale(d) + ")");

                    let extent = yscale[d].brush.extent();
                }

                // remove axis if dragged all the way left
                if (dragging[d] < 12 || dragging[d] > w-12) {
                    remove_axis(d,g);
                }

                // TODO required to avoid a bug
                xscale.domain(dimensions);
                update_ticks(d, extent);

                // rerender
                d3.select("#foreground").style("opacity", null);
                brush();
                delete this.__dragged__;
                delete this.__origin__;
                delete dragging[d];
            }));

    // Add an axis and title.
    g.append("svg:g")
        .attr("class", "axis")
        .attr("transform", "translate(0,0)")
        .each(function(d) { console.log("d", d); d3.select(this).call(axis.scale(yscale[d])); })
        .append("svg:text")
        .attr("text-anchor", "middle")
        .attr("y", function(d,i) { return i%2 == 0 ? -14 : -30 } )
        .attr("x", 0)
        .attr("class", "label")
        .text(String)
        .append("title")
        .text("Click to invert. Drag to reorder");

    // Add and store a brush for each axis.
    g.append("svg:g")
        .attr("class", "brush")
        .each(function(d) { d3.select(this).call(yscale[d].brush = d3.svg.brush().y(yscale[d]).on("brush", brush)); })
        .selectAll("rect")
        .style("visibility", null)
        .attr("x", -23)
        .attr("width", 36)
        .append("title")
        .text("Drag up or down to brush along this axis");

    g.selectAll(".extent")
        .append("title")
        .text("Drag or resize this filter");

    legend = create_legend(colors,brush);

    // Render full foreground
    brush();
});


///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// METHODS
///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// copy one canvas to another, grayscale
function gray_copy(source, target) {
    let pixels = source.getImageData(0,0,w,h);
    target.putImageData(grayscale(pixels),0,0);
}

// http://www.html5rocks.com/en/tutorials/canvas/imagefilters/
function grayscale(pixels, args) {
    let d = pixels.data;
    for (let i=0; i<d.length; i+=4) {
        let r = d[i];
        let g = d[i+1];
        let b = d[i+2];
        // CIE luminance for the RGB
        // The human eye is bad at seeing red and blue, so we de-emphasize them.
        let v = 0.2126*r + 0.7152*g + 0.0722*b;
        d[i] = d[i+1] = d[i+2] = v
    }
    return pixels;
};

/* ---------------------------------------------------
   Legend
   --------------------------------------------------- */

function create_legend(colors,brush) {
    // create legend
    let legend_data = d3.select("#legend")
        .html("")
        .selectAll(".row")
        .data( _.keys(colors).sort() )

    // filter by group
    let legend = legend_data
        .enter().append("div")
        .attr("title", "Hide group")
        .on("click", function(d) {
            // toggle food group
            if (_.contains(excluded_groups, d)) {
                d3.select(this).attr("title", "Hide group")
                excluded_groups = _.difference(excluded_groups,[d]);
                brush();
            } else {
                d3.select(this).attr("title", "Show group")
                excluded_groups.push(d);
                brush();
            }
        });

    legend
        .append("span")
        .style("background", function(d,i) { return color(d,0.85)})
        .attr("class", "color-bar");

    legend
        .append("span")
        .attr("class", "tally")
        .text(function(d,i) { return 0});

    legend
        .append("span")
        .text(function(d,i) { return " " + d});

    return legend;
}

/* ---------------------------------------------------
   Table
   --------------------------------------------------- */

// simple data table
function data_table(sample) {
    // sort by first column
    var sample = sample.sort(function(a,b) {
        let col = d3.keys(a)[0];
        return a[col] < b[col] ? -1 : 1;
    });

    let table = d3.select("#address-list")
        .html("")
        .selectAll(".row")
        .data(sample)
        .enter().append("div")
        .on("mouseover", highlight)
        .on("mouseout", unhighlight);

    table
        .append("span")
        .attr("class", "color-block")
        .style("background", function(d) { return color(d.cToken,0.85) })

    table
        .append("a")
        .attr("href", function(d) {
            return "https://ethstats.io/account/" + d.address;
        })
        .attr("target", function() {
            return "_blank";
        })
        .text(function(d) { return d.address; })
}

function search(selection,str) {
    pattern = new RegExp(str, "i")
    return _(selection).filter(function (d) {
        return pattern.exec(d.address);
    });
}


/* ---------------------------------------------------
   Rendering
   --------------------------------------------------- */

// render polylines i to i+render_speed
function render_range(selection, i, max, opacity) {
    selection.slice(i,max).forEach(function(d) {
        path(d, foreground, color(d.cToken,opacity));
    });
};

// Adjusts rendering speed
function optimize(timer) {
    let delta = (new Date()).getTime() - timer;
    render_speed = Math.max(Math.ceil(render_speed * 30 / delta), 8);
    render_speed = Math.min(render_speed, 300);
    return (new Date()).getTime();
}

// Feedback on rendering progress
function render_stats(i,n,render_speed) {
    d3.select("#rendered-count").text(i);
    d3.select("#rendered-bar")
        .style("width", (100*i/n) + "%");
    d3.select("#render-speed").text(render_speed);
}

// render a set of polylines on a canvas
function paths(selected, ctx, count) {
    let n = selected.length,
        i = 0,
        opacity = d3.min([2/Math.pow(n,0.25),1]),
        timer = (new Date()).getTime();

    selection_stats(opacity, n, data.length)

    shuffled_data = _.shuffle(selected);

    data_table(shuffled_data.slice(0,25));

    ctx.clearRect(0,0,w+1,h+1);

    // render all lines until finished or a new brush event
    function animloop(){
        if (i >= n || count < brush_count) return true;
        let max = d3.min([i+render_speed, n]);
        render_range(shuffled_data, i, max, opacity);
        render_stats(max,n,render_speed);
        i = max;
        timer = optimize(timer);  // adjusts render_speed
    };

    d3.timer(animloop);
}

// transition ticks for reordering, rescaling and inverting
function update_ticks(d, extent) {
    // update brushes
    if (d) {
        let brush_el = d3.selectAll(".brush")
            .filter(function(key) { return key == d; });
        // single tick
        if (extent) {
            // restore previous extent
            brush_el.call(yscale[d].brush = d3.svg.brush().y(yscale[d]).extent(extent).on("brush", brush));
        } else {
            brush_el.call(yscale[d].brush = d3.svg.brush().y(yscale[d]).on("brush", brush));
        }
    } else {
        // all ticks
        d3.selectAll(".brush")
            .each(function(d) { d3.select(this).call(yscale[d].brush = d3.svg.brush().y(yscale[d]).on("brush", brush)); })
    }

    brush_count++;

    show_ticks();

    // update axes
    d3.selectAll(".axis")
        .each(function(d,i) {
            // hide lines for better performance
            d3.select(this).selectAll('line').style("display", "none");

            // transition axis numbers
            d3.select(this)
                .transition()
                .duration(720)
                .call(axis.scale(yscale[d]));

            // bring lines back
            d3.select(this).selectAll('line').transition().delay(800).style("display", null);

            d3.select(this)
                .selectAll('text')
                .style('font-weight', null)
                .style('font-size', null)
                .style('display', null);
        });
}

/* ---------------------------------------------------
   Selecting
   --------------------------------------------------- */

// Handles a brush event, toggling the display of foreground lines.
// TODO refactor
function brush() {
    brush_count++;
    let actives = dimensions.filter(function(p) { return !yscale[p].brush.empty(); }),
        extents = actives.map(function(p) { return yscale[p].brush.extent(); });

    // hack to hide ticks beyond extent
    let b = d3.selectAll('.dimension')[0]
        .forEach(function(element, i) {
            let dimension = d3.select(element).data()[0];
            if (_.include(actives, dimension)) {
                let extent = extents[actives.indexOf(dimension)];
                d3.select(element)
                    .selectAll('text')
                    .style('font-weight', 'bold')
                    .style('font-size', '13px')
                    .style('display', function() {
                        let value = d3.select(this).data();
                        return extent[0] <= value && value <= extent[1] ? null : "none"
                    });
            } else {
                d3.select(element)
                    .selectAll('text')
                    .style('font-size', null)
                    .style('font-weight', null)
                    .style('display', null);
            }
            d3.select(element)
                .selectAll('.label')
                .style('display', null);
        });
    ;

    // bold dimensions with label
    d3.selectAll('.label')
        .style("font-weight", function(dimension) {
            if (_.include(actives, dimension)) return "bold";
            return null;
        });

    // Get lines within extents
    let selected = [];
    data
        .filter(function(d) {
            return !_.contains(excluded_groups, d.cToken);
        })
        .map(function(d) {
            return actives.every(function(p, dimension) {
                return extents[dimension][0] <= d[p] && d[p] <= extents[dimension][1];
            }) ? selected.push(d) : null;
        });

    // free text search
    let query = d3.select("#search")[0][0].value;
    if (query.length > 0) {
        selected = search(selected, query);
    }

    if (selected.length < data.length && selected.length > 0) {
        d3.select("#keep-data").attr("disabled", null);
        d3.select("#exclude-data").attr("disabled", null);
    } else {
        d3.select("#keep-data").attr("disabled", "disabled");
        d3.select("#exclude-data").attr("disabled", "disabled");
    };

    // total by food group
    let tallies = _(selected)
        .groupBy(function(d) { return d.cToken; })

    // include empty groups
    _(colors).each(function(v,k) { tallies[k] = tallies[k] || []; });

    legend
        .style("text-decoration", function(d) { return _.contains(excluded_groups,d) ? "line-through" : null; })
        .attr("class", function(d) {
            return (tallies[d].length > 0)
                ? "row"
                : "row off";
        });

    legend.selectAll(".color-bar")
        .style("width", function(d) {
            return Math.ceil(600*tallies[d].length/data.length) + "px"
        });

    legend.selectAll(".tally")
        .text(function(d,i) { return tallies[d].length });

    // Render selected lines
    paths(selected, foreground, brush_count, true);
}

// Feedback on selection
function selection_stats(opacity, n, total) {
    d3.select("#data-count").text(total);
    d3.select("#selected-count").text(n);
    d3.select("#selected-bar").style("width", (100*n/total) + "%");
    d3.select("#opacity").text((""+(opacity*100)).slice(0,4) + "%");
}

// Highlight single polyline
function highlight(d) {
    d3.select("#foreground").style("opacity", "0.25");
    d3.selectAll(".row").style("opacity", function(p) { return (d.cToken == p) ? null : "0.3" });
    path(d, highlighted, color(d.cToken,1));
}

// Remove highlight
function unhighlight() {
    d3.select("#foreground").style("opacity", null);
    d3.selectAll(".row").style("opacity", null);
    highlighted.clearRect(0,0,w,h);
}

function position(d) {
    let v = dragging[d];
    return v == null ? xscale(d) : v;
}

function invert_axis(d) {
    // save extent before inverting
    if (!yscale[d].brush.empty()) {
        let extent = yscale[d].brush.extent();
    }
    if (yscale[d].inverted == true) {
        yscale[d].range([h, 0]);
        d3.selectAll('.label')
            .filter(function(p) { return p == d; })
            .style("text-decoration", null);
        yscale[d].inverted = false;
    } else {
        yscale[d].range([0, h]);
        d3.selectAll('.label')
            .filter(function(p) { return p == d; })
            .style("text-decoration", "underline");
        yscale[d].inverted = true;
    }
    return extent;
}

function path(d, ctx, color) {
    if (color) ctx.strokeStyle = color;
    ctx.beginPath();
    let x0 = xscale(0)-15,
        y0 = yscale[dimensions[0]](d[dimensions[0]]);   // left edge
    ctx.moveTo(x0,y0);
    dimensions.map(function(p,i) {
        let x = xscale(p),
            y = yscale[p](d[p]);
        let cp1x = x - 0.88*(x-x0);
        let cp1y = y0;
        let cp2x = x - 0.12*(x-x0);
        let cp2y = y;
        ctx.bezierCurveTo(cp1x, cp1y, cp2x, cp2y, x, y);
        x0 = x;
        y0 = y;
    });
    ctx.lineTo(x0+15, y0);                               // right edge
    ctx.stroke();
};

function color(d,a) {
    // console.log("d", d);
    let c = colors[d];
    return ["hsla(",c[0],",",c[1],"%,",c[2],"%,",a,")"].join("");
}

// Rescale to new dataset domain
function rescale() {
    // reset yscales, preserving inverted state
    dimensions.forEach(function(d,i) {
        if (yscale[d].inverted) {
            yscale[d] = d3.scale.log()
                .domain(d3.extent(data, function(p) { return +p[d]; }))
                .range([0, h]);
            yscale[d].inverted = true;
        } else {
            yscale[d] = d3.scale.log()
                .domain(d3.extent(data, function(p) { return +p[d]; }))
                .range([h, 0]);
        }
    });

    update_ticks();

    // Render selected data
    paths(data, foreground, brush_count);
}

// Get polylines within extents
function actives() {
    let actives = dimensions.filter(function(p) { return !yscale[p].brush.empty(); }),
        extents = actives.map(function(p) { return yscale[p].brush.extent(); });

    // filter extents and excluded groups
    let selected = [];
    data
        .filter(function(d) {
            return !_.contains(excluded_groups, d.cToken);
        })
        .map(function(d) {
            return actives.every(function(p, i) {
                return extents[i][0] <= d[p] && d[p] <= extents[i][1];
            }) ? selected.push(d) : null;
        });

    // free text search
    let query = d3.select("#search")[0][0].value;
    if (query > 0) {
        selected = search(selected, query);
    }

    return selected;
}

// scale to window size
window.onresize = function() {
    width = document.body.clientWidth,
        height = d3.max([document.body.clientHeight-600, 220]);

    w = width - m[1] - m[3],
        h = height - m[0] - m[2];

    d3.select("#chart")
        .style("height", (h + m[0] + m[2]) + "px")

    d3.selectAll("canvas")
        .attr("width", w)
        .attr("height", h)
        .style("padding", m.join("px ") + "px");

    d3.select("svg")
        .attr("width", w + m[1] + m[3])
        .attr("height", h + m[0] + m[2])
        .select("g")
        .attr("transform", "translate(" + m[3] + "," + m[0] + ")");

    xscale = d3.scale.ordinal().rangePoints([0, w], 1).domain(dimensions);
    dimensions.forEach(function(d) {
        // TODO . confditional for this
        // console.log("d");
        if (d !== "type") {
            yscale[d].range([h, 0]);
        }
    });

    d3.selectAll(".dimension")
        .attr("transform", function(d) { return "translate(" + xscale(d) + ")"; })
    // update brush placement
    d3.selectAll(".brush")
        .each(function(d) { d3.select(this).call(yscale[d].brush = d3.svg.brush().y(yscale[d]).on("brush", brush)); })
    brush_count++;

    // update axis placement
    axis = axis.ticks(1+height/50),
        d3.selectAll(".axis")
            .each(function(d) { d3.select(this).call(axis.scale(yscale[d])); });

    // render data
    brush();
};


/* ---------------------------------------------------
   Data
   --------------------------------------------------- */

// Remove all but selected from the dataset
function keep_data() {
    new_data = actives();
    if (new_data.length == 0) {
        alert("I don't mean to be rude, but I can't let you remove all the data.\n\nTry removing some brushes to get your data back. Then click 'Keep' when you've selected data you want to look closer at.");
        return false;
    }
    data = new_data;
    rescale();
}

// Exclude selected from the dataset
function exclude_data() {
    new_data = _.difference(data, actives());
    if (new_data.length == 0) {
        alert("I don't mean to be rude, but I can't let you remove all the data.\n\nTry selecting just a few data points then clicking 'Exclude'.");
        return false;
    }
    data = new_data;
    rescale();
}

// Export data
function export_csv() {
    let keys = d3.keys(data[0]);
    let rows = actives().map(function(row) {
        return keys.map(function(k) { return row[k]; })
    });
    let csv = d3.csv.format([keys].concat(rows)).replace(/\n/g,"<br/>\n");
    let styles = "<style>body { font-family: sans-serif; font-size: 12px; }</style>";
    window.open("text/csv").document.write(styles + csv);
}

function remove_axis(d,g) {
    dimensions = _.difference(dimensions, [d]);
    xscale.domain(dimensions);
    g.attr("transform", function(p) { return "translate(" + position(p) + ")"; });
    g.filter(function(p) { return p == d; }).remove();
    update_ticks();
}

d3.select("#keep-data").on("click", keep_data);
d3.select("#exclude-data").on("click", exclude_data);
d3.select("#export-data").on("click", export_csv);
d3.select("#search").on("keyup", brush);


/* ---------------------------------------------------
   Appearance toggles
   --------------------------------------------------- */

d3.select("#hide-ticks").on("click", hide_ticks);
d3.select("#show-ticks").on("click", show_ticks);
d3.select("#dark-theme").on("click", dark_theme);
d3.select("#light-theme").on("click", light_theme);

function hide_ticks() {
    d3.selectAll(".axis g").style("display", "none");
    // d3.selectAll(".axis path").style("display", "none");
    d3.selectAll(".background").style("visibility", "hidden");
    d3.selectAll("#hide-ticks").attr("disabled", "disabled");
    d3.selectAll("#show-ticks").attr("disabled", null);
};

function show_ticks() {
    d3.selectAll(".axis g").style("display", null);
    //d3.selectAll(".axis path").style("display", null);
    d3.selectAll(".background").style("visibility", null);
    d3.selectAll("#show-ticks").attr("disabled", "disabled");
    d3.selectAll("#hide-ticks").attr("disabled", null);
};

function dark_theme() {
    d3.select("body").attr("class", "dark");
    d3.selectAll("#dark-theme").attr("disabled", "disabled");
    d3.selectAll("#light-theme").attr("disabled", null);
}

function light_theme() {
    d3.select("body").attr("class", null);
    d3.selectAll("#light-theme").attr("disabled", "disabled");
    d3.selectAll("#dark-theme").attr("disabled", null);
}

function wrap(text, width) {
    text.each(function() {
        let text = d3.select(this),
            words = text.text().split(/\s+/).reverse(),
            word,
            line = [],
            lineNumber = 0,
            lineHeight = 1.1, // ems
            y = text.attr("y"),
            dy = parseFloat(text.attr("dy")),
            tspan = text.text(null).append("tspan").attr("x", 0).attr("y", y).attr("dy", dy + "em");
        while (word = words.pop()) {
            line.push(word);
            tspan.text(line.join(" "));
            if (tspan.node().getComputedTextLength() > width) {
                line.pop();
                tspan.text(line.join(" "));
                line = [word];
                tspan = text.append("tspan").attr("x", 0).attr("y", y).attr("dy", ++lineNumber * lineHeight + dy + "em").text(word);
            }
        }
    });
}