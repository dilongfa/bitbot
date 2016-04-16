var React = require('react'),
    ReactDOM = require('react-dom'),
    hashHistory = require('react-router').hashHistory;

var SelectField = require('material-ui/lib/select-field'),
    MenuItem = require('material-ui/lib/menus/menu-item');

var pairs = require('./pairs');

var CreateLineChart = require('./line-chart');

module.exports = React.createClass({
    getInitialState: function () {
        return {data: []};
    },

    componentDidMount: function () {
        this._updateState(this.props.location);
    },

    componentWillReceiveProps: function (nextProps) {
        this._updateState(nextProps.location)
    },

    _updateState: function (location) {
        var that = this;
        $.get(location.pathname, location.query, function (data) {
            that.setState({data: data});
        });
    },

    render: function () {
        return <div>
            <h1>Bid/Ask</h1>
            <SearchForm location={this.props.location} pair={this.props.params.pair} />
            <Table data={this.state.data} />
        </div>
    }
});

var SearchForm = React.createClass({
    handleChange: function (e, i, pair) {
        this._submit(pair);
        e.preventDefault();
    },

    handleSubmit: function (e) {
        this._submit(this.props.pair);
        e.preventDefault();
    },

    _submit: function (pair) {
        var form = ReactDOM.findDOMNode(this);
        hashHistory.push('/bid_ask/' + pair);
    },

    render: function () {
        return <form onSubmit={this.handleSubmit}>
            <SelectField value={this.props.pair} onChange={this.handleChange}>
                {pairs.map(function (p) {
                    return <MenuItem value={p.symbol} primaryText={p.label} />
                })}
            </SelectField>
            {/* TODO: onSubmit isn't triggered whithout if the form doesn't contain that button.
            I don't understand why... */}
            <input type="submit" value="send" />
        </form>
    }
});

var Table =  React.createClass({
    render: function () {
        if (this.props.data.length == 0)
            return <div />

        var last = this.props.data[this.props.data.length - 1] || {};
        // TODO: share this list
        var exchangers = ['Cex', 'Kraken', 'Btce', 'Hitbtc', 'Bitfinex'];
        var data = [];

        for (var i=0; i < exchangers.length; i++) {
            var exchanger = exchangers[i];
            if (last.Orderbooks[exchanger]) {
                data.push(last.Orderbooks[exchanger])
            }
        }

        var props = this.props;

        var rows = data.map(function (r, i) {
            return <tr>
                <td>{r.Exchanger}</td>
                <td>{r.Bids[0].Price}</td>
                <td>{r.Asks[0].Price}</td>
                <td ref={function (el) {
                    var chart = CreateLineChart(props.data, r.Exchanger);
                    if (el != null) {
                        el.innerHTML = "";
                        el.appendChild(chart);
                    }
                }}></td>
            </tr>
        });

        return <table>
            <thead>
                <tr>
                    <th>Exchanger</th>
                    <th>Bid</th>
                    <th>Ask</th>
                    <th>Evol</th>
                </tr>
            </thead>
            <tbody>
                {rows}
            </tbody>
        </table>
    }
});
