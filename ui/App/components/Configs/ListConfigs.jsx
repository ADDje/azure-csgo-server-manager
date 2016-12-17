import React from 'react';
import swal from 'sweetalert';

class ListConfigs extends React.Component {

    constructor(params) {
        super(params)

        this.changeConfig = this.changeConfig.bind(this)
        
        this.createConfig = this.createConfig.bind(this)
        this.createClick = this.createClick.bind(this)

        this.deleteClick = this.deleteClick.bind(this)
        this.deleteConfig = this.deleteConfig.bind(this)
    }

    changeConfig(config, configName, e) {
        e.preventDefault()

        this.props.focusConfig(config, configName)
    }

    createConfig(name) {           
        $.post({
            url: "api/configs/create/" + name,
            dataType: "json",
            success: (resp) => {

                swal({
                    timer: 2000,
                    title: "Nice!",
                    text: name + " created!",
                    type: "success"
                })
                
                this.props.reloadConfigs();
            }
        })
    }

    deleteClick(name, e) {
        e.preventDefault();
        swal({
            title: "Are you sure?",
            text: "You will not be able to recover " + name,
            type: "warning",
            showCancelButton: true,
            confirmButtonColor: "#DD6B55",
            confirmButtonText: "Yes, delete it!",
            closeOnConfirm: false,
            showLoaderOnConfirm: true
        },
        function(){
            this.deleteConfig(name)
        }.bind(this));
    }

    deleteConfig(name) {
        $.post({
            url: "api/configs/delete/" + name,
            dataType: "json",
            success: (resp) => {

                swal({
                    timer: 2000,
                    title: "Nice!",
                    text: name + " deleted!",
                    type: "success"
                })
                
                this.props.reloadConfigs();
            }
        })
    }

    createClick(e) {
        e.preventDefault()
        swal({
                title: "Create config",
                text: "Name your new config. Must end in .cfg and not already exist.",
                type: "input",
                showCancelButton: true,
                closeOnConfirm: false,
                animation: "slide-from-top",
                inputPlaceholder: "myconfig.cfg",
                showLoaderOnConfirm: true
            },
            function(inputValue) {
                if (inputValue === false) return false;
                
                inputValue = inputValue.trim()

                if (inputValue === "") {
                    swal.showInputError("You need to write something!")
                    return false
                }

                if (inputValue.substr(-4) !== ".cfg") {
                    swal.showInputError("Must end in .cfg!")
                    return false
                }

                if (inputValue.length < 5 || inputValue.indexOf(" ") !== -1) {
                    swal.showInputError("Invalid filename!")
                    return false
                }

                this.createConfig(inputValue)
            }.bind(this))
    }

    render() {
        return(
            <div className="box">
                <div className="box-header">
                    <h3 className="box-title">Manage Server Configs</h3>
                </div>
                
                <div className="box-body">
                    <div className="table-responsive">
                        <table className="table table-striped">
                            <thead>
                                <tr>
                                    <th>Name</th>
                                </tr>
                            </thead>
                            <tbody>
                            {Object.keys(this.props.configs).map ( (config, i) => {
                                return(
                                    <tr key={i}>
                                        <td>
                                            <a className="row-link"
                                                href="#"
                                                onClick={this.changeConfig.bind(this, this.props.configs[config], config)}
                                            >
                                                {config}
                                                <button className="btn btn-sm btn-danger" onClick={this.deleteClick.bind(this, config)}>
                                                    <i className="fa fa-trash" />
                                                </button>
                                            </a>
                                        </td>
                                    </tr>
                                )                                       
                            })}
                            </tbody>
                        </table>
                        <button className="btn btn-sm btn-success" onClick={this.createClick}>
                            <i className="fa fa-plus" />
                        </button>
                    </div>
                </div>
            </div>
        )
    }
}

ListConfigs.propTypes = {
    configs: React.PropTypes.object.isRequired,
    focusConfig: React.PropTypes.func.isRequired,
    reloadConfigs: React.PropTypes.func.isRequired
}

export default ListConfigs
