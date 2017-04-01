import React, {Component, PropTypes} from 'react'
import {connect} from 'react-redux'
import {Media, Row, Col, Grid, Panel, ListGroup, ListGroupItem} from 'react-bootstrap'
import {fetchProfile, addTweet} from '../actions/twitterActions'

class Profile extends Component {

    constructor(props) {
        super(props)
        this.state = {text: ''}

        this.handleChange = this.handleChange.bind(this)
        this.handleSubmit = this.handleSubmit.bind(this)
    }

    componentDidMount() {
        this.props.dispatch(fetchProfile())
    }

    handleChange(event) {
        this.setState({text: event.target.value})
    }

    handleSubmit(event) {
        const text = this.state.text.trim()
        this.props.dispatch(addTweet(text))
        this.setState({text: ''})
        event.preventDefault()
    }

    render() {
        return (
            <Panel>
                <ListGroup fill>
                    <ListGroupItem>
                        <Media>
                            <Media.Body>
                                 <Media.Heading>{this.props.user.name}</Media.Heading>
                            </Media.Body>
                        </Media>
                    </ListGroupItem>
                    <ListGroupItem>
                        <Grid fluid>
                            <Row className="show-grid">
                                <Col sm={4} md={4}>Tweets<span className="badge">{this.props.user.tweets}</span></Col>
                                <Col sm={4} md={4}>Following<span className="badge">{this.props.user.following}</span></Col>
                                <Col sm={4} md={4}>Followers<span className="badge">{this.props.user.followers}</span></Col>
                            </Row>
                        </Grid>
                    </ListGroupItem>
                </ListGroup>
                <form onSubmit={this.handleSubmit}>
                    <input className="form-control" type="text" placeholder="Compose new tweet..." value={this.state.text} onChange={this.handleChange}/>
                    <input type="submit" value="Submit" />
                </form>
            </Panel>
        )
    }
}

Profile.propTypes = {
    user: PropTypes.shape({
        name: PropTypes.string.isRequired
    }).isRequired
}

const mapStateToProps = (state) => {
    return {
        user: state.user
    }
}

export default connect(mapStateToProps)(Profile)