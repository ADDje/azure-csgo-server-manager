import React from 'react';

class Settings extends React.Component {
    constructor(props) {
        super(props)
    }

    render() {
        return(
            <tbody>
            {Object.keys(this.props.config).map(function(key) {
                return(
                    <tr key={key}>
                        <td>{key}</td>
                        <td>{this.props.config[key]}</td>
                    </tr>
                )                                                  
            }, this)}        
            </tbody>
        )
    }

}

Settings.propTypes = {
    config: React.PropTypes.object.isRequired,
    section: React.PropTypes.string.isRequired,
}

export default Settings
