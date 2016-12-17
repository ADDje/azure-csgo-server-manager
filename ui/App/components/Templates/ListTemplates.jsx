import React from 'react';

class ListTemplates extends React.Component {

    constructor(params) {
        super(params)

        this.changeTemplate = this.changeTemplate.bind(this)

        this.createTemplate = this.createTemplate.bind(this)
        this.createClick = this.createClick.bind(this)

        this.deleteClick = this.deleteClick.bind(this)
        this.deleteTemplate = this.deleteTemplate.bind(this)
    }

    changeTemplate(template, templateName, e) {
        e.preventDefault()

        this.props.focusTemplate(template, templateName)
    }

    createTemplate(name) {           
        $.post({
            url: "api/templates/create/" + name,
            dataType: "json",
            success: (resp) => {

                swal({
                    timer: 2000,
                    title: "Nice!",
                    text: name + " created!",
                    type: "success"
                })
                
                this.props.reloadTemplates();
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
            this.deleteTemplate(name)
        }.bind(this));
    }

    deleteTemplate(name) {
        $.post({
            url: "api/templates/delete/" + name,
            dataType: "json",
            success: (resp) => {

                swal({
                    timer: 2000,
                    title: "Nice!",
                    text: name + " deleted!",
                    type: "success"
                })
                
                this.props.reloadTemplates();
            }
        })
    }

    createClick(e) {
        e.preventDefault()
        swal({
                title: "Create template",
                text: "Name your new template. Do not include .json as this will be added automatically.",
                type: "input",
                showCancelButton: true,
                closeOnConfirm: false,
                animation: "slide-from-top",
                inputPlaceholder: "myTemplate",
                showLoaderOnConfirm: true
            },
            function(inputValue) {
                if (inputValue === false) return false;
                
                inputValue = inputValue.trim()

                if (inputValue === "") {
                    swal.showInputError("You need to write something!")
                    return false
                }

                if (inputValue.indexOf(" ") !== -1 || inputValue.indexOf(".") !== -1) {
                    swal.showInputError("Invalid filename!")
                    return false
                }

                this.createTemplate(inputValue)
            }.bind(this))
    }

    render() {
        return(
            <div className="box">
                <div className="box-header">
                    <h3 className="box-title">Manage Deployment Templates</h3>
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
                            {Object.keys(this.props.templates).map ( (template, i) => {
                                return(
                                    <tr key={i}>
                                        <td>
                                            <a className="row-link"
                                                href="#"
                                                onClick={this.changeTemplate.bind(this, this.props.templates[template], template)}
                                            >
                                                {template}
                                                <button className="btn btn-sm btn-danger" onClick={this.deleteClick.bind(this, template)}>
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

ListTemplates.propTypes = {
    focusTemplate: React.PropTypes.func.isRequired,
    reloadTemplates: React.PropTypes.func.isRequired,
    templates: React.PropTypes.object.isRequired,
}

export default ListTemplates
