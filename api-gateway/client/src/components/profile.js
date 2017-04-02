import React, {Component, PropTypes} from 'react'
import {connect} from 'react-redux'
import {Media, Row, Col, Grid, Panel, ListGroup, ListGroupItem} from 'react-bootstrap'
import {fetchProfile} from '../actions/twitterActions'
import TweetInput from './tweetInput'

class Profile extends Component {

    componentDidMount() {
        this.props.dispatch(fetchProfile())
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
                <TweetInput/>
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