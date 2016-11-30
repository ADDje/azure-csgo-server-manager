import React from 'react';
import update from 'immutability-helper';

class DynamicConfig extends React.Component {

    constructor(props) {
        super(props);

        this.state = {
            config: null,
            configName: this.props.configName,
            configParameters: {}
        }

        this.changeField = this.changeField.bind(this)
    }

    componentWillReceiveProps(nextProps) {
        if(nextProps.configName !== this.props.configName &&
            nextProps.config !== null) {
            // Populate the parameters
            Object.assign(this.state.configParameters, nextProps.config)
        }
    }

    changeField(event) {
        console.log(event)
    }

    increaseField(key) {
        var param = {}
        param[key] = {$set: this.state.configParameters[key] + 1};
        this.setState({
            configParameters: update(this.state.configParameters, param)
        })
    }

    decreaseField(key) {
        var param = {}
        param[key] = {$set: this.state.configParameters[key] - 1};
        this.setState({
            configParameters: update(this.state.configParameters, param)
        })
    }

    render() {
        if(this.props.config === null) {
            return null;
        }

        var fields = []
        for(var key in this.props.config) {
            var buttons = null;
            if(typeof(this.props.config[key]) !== "string") {
                buttons = (<div className="input-group-btn">
                        <button type="button" className="btn btn-primary" onClick={this.increaseField.bind(this, key)}>
                            <i className="fa fa-arrow-up" />
                        </button>
                        <button type="button" className="btn btn-primary" onClick={this.decreaseField.bind(this, key)}>
                            <i className="fa fa-arrow-down" />
                        </button>
                    </div>)
            }
 
            fields.push(<div key={key} className="dynamic-config-field">
                <label htmlFor={key}>{key}</label>
                <div id={key} className="input-group">
                    <input ref={key} name={key} id={key} type="text" className="form-control" onChange={this.changeField} value={this.state.configParameters[key]} />
                    {buttons}
                </div>
            </div>)
        }

        return (
            <div className="dynamic-config">
                {fields}
            </div>
        );
    }

}

export default DynamicConfig