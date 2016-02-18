var React = require('react');
var d3 = require('d3');

var WIDTH = 600,
    HEIGHT = 100,
    margin = {top: 20, right: 20, bottom: 30, left: 50};

var formatDate = d3.time.format("%Y-%m-%d %H:%M");

// TODO: share this list
var exchangers = ['Cex', 'Kraken', 'Btce', 'Hitbtc', 'Bitfinex'];

module.exports = React.createClass({
    componentDidMount: function() {
        this._render(this.props);
    },

    componentWillReceiveProps: function(props, state) {
        this._render(props);
    },

    _render: function (props) {
        var el = this.getDOMNode();
        el.innerHTML = "";
        el.appendChild(createLineCharts(el, props.data));
    },

    render: function() {
        return <div className="line-charts" />
    }
});

function createLineCharts(el, data) {
    var dates = [];
    for (var i=0; i < data.length; i++) {
        dates.push(formatDate.parse(data[i].StartDate));
    }

    var svgRoot = document.createElementNS(d3.ns.prefix.svg, 'svg');

    var svg = d3.select(svgRoot)
        .attr("width", WIDTH)
        .attr("height", HEIGHT * exchangers.length);

    for (var i = 0; i < exchangers.length; i++) {
        var rect = document.createElementNS(d3.ns.prefix.svg, 'rect');

        d3.select(rect)
            .attr("width", WIDTH)
            .attr("height", HEIGHT)
            .attr("transform", "translate(0," + HEIGHT*i + ")");

        svgRoot.appendChild(rect);
    }

    for (var i = 0; i < exchangers.length; i++) {
        var g = addLineChart(dates, data, exchangers[i]);
        d3.select(g).attr("transform", "translate(0," + HEIGHT*i + ")");
        svgRoot.appendChild(g);
    }

    var layer = addInteractiveLayer(dates.length, svg);
    svgRoot.appendChild(layer);
    return svgRoot;
}

function addLineChart(dates, data, exchanger) {
    var width = WIDTH - margin.left - margin.right,
        height = HEIGHT - margin.top - margin.bottom;

    // x axis
    var x = d3.time.scale().range([0, width]);
    x.domain(d3.extent(dates));

    var xAxis = d3.svg
        .axis()
        .scale(x)
        .orient("bottom")
        .ticks(d3.time.minutes, 10)
        .tickFormat(d3.time.format("%H:%M"));

    // y axis
    var ymin = d3.min(data, function(d) { return d.Orderbooks[exchanger] ? d.Orderbooks[exchanger].Bids[0].Price: null }),
        ymax = d3.max(data, function(d) { return d.Orderbooks[exchanger] ? d.Orderbooks[exchanger].Asks[0].Price: null }),
        delta = (ymax - ymin) * 0.1;

    var y = d3.scale.linear()
        .domain([ymin - delta, ymax + delta])
        .range([height, 0]);

    var yAxis = d3.svg
        .axis()
        .scale(y)
        .orient("left")
        .ticks(3);

    var g = document.createElementNS(d3.ns.prefix.svg, 'g');

    var container = d3.select(g)
        .attr("class", "chart")
        .attr("width", width)
        .attr("height", height)
        .append("g")
        .attr("transform", "translate(" + margin.left + "," + margin.top + ")");

    container.append("g")
        .attr("class", "x axis")
        .attr("transform", "translate(0," + height + ")")
        .call(xAxis);

    container.append("g")
        .attr("class", "y axis")
        .call(yAxis)
        .append("text")
        .attr("transform", "rotate(-90)")
        .attr("y", 6)
        .attr("dy", ".71em")
        .style("text-anchor", "end");
        // .text("Price ($)");

    addLine(container, x, y, dates, data, exchanger, 'Bids', 'steelblue');
    addLine(container, x, y, dates, data, exchanger, 'Asks', '#FC9E27');
    return g;
}

// TODO: data are ordered by ts desc
function addLine(container, x, y, dates, data, exchanger, attr, color) {
    var segmentData = [];

    for (var i = 0; i < data.length; i++) {
        if (data[i].Orderbooks[exchanger]) {
            segmentData.push({
                date: dates[i],
                value: data[i].Orderbooks[exchanger][attr][0].Price,
            });
        } else if (segmentData.length > 0) {
            addSegment();
            segmentData = [];
        }
    }

    if (segmentData.length > 0) {
        addSegment();
    }

    function addSegment() {
        var line = d3.svg.line()
            .x(function(d) { return x(d.date); })
            .y(function(d) { return y(d.value); });

        container.append("path")
            .datum(segmentData)
            .attr("class", "line")
            .attr("d", line)
            .attr("stroke", color);
    }
}

function addInteractiveLayer(n) {
    var g = document.createElementNS(d3.ns.prefix.svg, 'g');
    var innerWidth = WIDTH - margin.left - margin.right;

    var wrapper = d3.select(g)
        .attr('class', 'interactive-wrap');

    var line = wrapper.append('line')
        .attr('y1', 0)
        .attr('y2', HEIGHT*5)
        .attr('x1', margin.left)
        .attr('x2', margin.left)
        .attr('dx', 1)
        .attr('stroke', '#ccc')
        .attr('stroke-width', 1);

    var rect = wrapper.append('rect')
        .attr('width', innerWidth)
        .attr('height', HEIGHT*5)
        .attr('opacity', '0')
        .attr('pointer-events', 'all')
        .attr('transform', 'translate(' + margin.left + ',0)')
        .on('mouseover', mouseover)
        .on('mousemove', mousemove)
        .on('mouseout', mouseout);

    function mouseover() {
        line.attr('opacity', 1);
    }

    var dt = innerWidth / n;

    function mousemove() {
        var xy = d3.mouse(g);
        var idx = Math.floor((xy[0] - margin.left) / dt);
        // emit custom event here/ ?
        line.attr('x1', xy[0]);
        line.attr('x2', xy[0]);
    }

    function mouseout() {
        line.attr('opacity', 0);
    }

    return g;
}
