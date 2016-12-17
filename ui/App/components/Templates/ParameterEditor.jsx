import React from 'react';
import update from 'immutability-helper';

class ParameterEditor extends React.Component {

    constructor(props) {
        super(props)

        var parameters = null
        if (this.props.parameters !== undefined) {
            parameters = this.props.parameters.parameters
        }

        this.state = {
            parameters: parameters,
            templateName: this.props.templateName,
            templateParameters: {},
            deleted: []
        }

        if (this.props.parameters !== undefined) { 
            Object.assign(this.state.templateParameters, this.props.parameters)
        }

        this.changeField = this.changeField.bind(this)
        this.save = this.save.bind(this)

        this.createClick = this.createClick.bind(this)
        this.deleteClick = this.deleteClick.bind(this)
        this.deleteParameter = this.deleteParameter.bind(this)
    }

    componentWillReceiveProps(nextProps) {
        if (nextProps.templateName !== this.props.templateName &&
            nextProps.parameters !== null) {

            var newParameters = {}
            Object.assign(newParameters, nextProps.parameters)
            this.setState({
                templateParameters: newParameters,
                deleted: []
            })
        }
    }

    changeField(key, event) {
        var param = {}
        param[key] = {value: {$set: event.target.value}}
        this.setState({
            templateParameters: update(this.state.templateParameters, param)
        })
    }

    increaseField(key) {
        var param = {}
        param[key] = {value: {$set: this.state.configParameters[key] + 1}}
        this.setState({
            templateParameters: update(this.state.configParameters, param)
        })
    }

    decreaseField(key) {
        var param = {}
        param[key] = {value: {$set: this.state.configParameters[key] - 1}}
        this.setState({
            templateParameters: update(this.state.configParameters, param)
        })
    }

    getValue(val) {
        if (val === null) {
            return ""
        }
        return val
    }

    createClick(e) {
        e.preventDefault()
        swal({
                title: "Add parameters",
                text: "Parameter Name:",
                type: "input",
                showCancelButton: true,
                animation: "slide-from-top",
                inputPlaceholder: "someParameter"
            },
            function(inputValue) {
                if (inputValue === false) return false;
                
                inputValue = inputValue.trim()

                if (inputValue === "") {
                    swal.showInputError("You need to write something!")
                    return false
                }

                if (inputValue.indexOf(" ") !== -1) {
                    swal.showInputError("Invalid parameter!")
                    return false
                }

                var newProp = {};
                newProp[inputValue] = { value: "" };
                this.setState({
                    templateParameters: update(this.state.templateParameters, {$merge: newProp})
                })
            }.bind(this))
    }

    deleteClick(name, e) {
        e.preventDefault();
        swal({
            title: "Are you sure?",
            text: "You will not be able to recover parameter: " + name,
            type: "warning",
            showCancelButton: true,
            confirmButtonColor: "#DD6B55",
            confirmButtonText: "Yes, delete it!",
            closeOnConfirm: false,
            showLoaderOnConfirm: true
        },
        function(){
            this.deleteParameter(name)
        }.bind(this));
    }

    deleteParameter(name) {
        var params = {}
        Object.assign(params, this.state.templateParameters)
        delete params[name]
        this.setState({
            templateParameters: params,
            deleted: update(this.state.deleted, {$push: [name]})
        })

        swal({
            timer: 100,
            title: "Nice!",
            text: name + " deleted!",
            type: "success"
        })
                
    }

    save() {
        // Start out with the original
        var newContent = this.props.parameters

        // Replace parameters with new ones
        for (var x in this.state.templateParameters) {
            newContent[x] = this.state.templateParameters[x]
        }

        // Delete any that shouldn't be there
        for (var y in this.state.deleted) {
            delete newContent[this.state.deleted[y]]
        }

        $.ajax({
            type: "POST",
            url: "/api/templates/" + this.state.templateName + "/parameters",
            dataType: "json",
            data: JSON.stringify(newContent, null, 4),
            success: (resp) => {
                if (typeof(resp.success) === "undefined" || resp.success === false) {
                    this.setState({isLoading: false, error: resp.data});
                } else {
                    this.setState({isLoading: false, error: null});
                    this.props.reloadSelected()
                }
            }
        })
    }

    render() {
        if (this.props.parameters === null || this.props.parameters === undefined) {
            return null;
        }

        var fields = []
        var parameters = this.state.templateParameters
        for (var key in parameters) {
            var buttons = null;
            if (typeof(parameters[key].value) !== "string") {
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
                <button className="btn btn-sm btn-danger" onClick={this.deleteClick.bind(this, key)}>
                    <i className="fa fa-trash" />
                </button>
                <div id={key} className="input-group">
                    <input ref={key} name={key} id={key} type="text" className="form-control" onChange={this.changeField.bind(this, key)} value={this.getValue(this.state.templateParameters[key].value)} />
                    {buttons}
                </div>
            </div>)
        }

        return (
            <div className="template-parameters">
                <div className="lead">
                    Parameters also support variables. "{"${n}"}" is the vm number
                </div>
                {fields}
                <button className="btn btn-sm btn-success" onClick={this.createClick}>
                    <i className="fa fa-plus" />
                </button>
                <div style={{marginTop: 10 + 'px'}}>
                    <button onClick={this.save} type="submit" className="btn btn-primary">Submit</button>
                </div>
            </div>
        );
    }

}

ParameterEditor.propTypes = {
    reloadSelected: React.PropTypes.func.isRequired,
}

export default ParameterEditor