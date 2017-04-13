import {combineReducers} from 'redux'
import {
    LOGIN_REQUEST, LOGIN_SUCCESS, LOGIN_FAILURE, LOGOUT_SUCCESS, RECEIVE_TWEETS, RECEIVE_PROFILE, RECEIVE_STATS
} from '../actions/twitterActions'

function tweets(state = [], action) {
    switch (action.type) {
        case RECEIVE_TWEETS:
            return Object.assign([], state, action.tweets)
        default:
            return state
    }
}

function user(state = {
                  name: 'loading...'
              }, action) {
    switch (action.type) {
        case RECEIVE_PROFILE:
            return Object.assign({}, state, action.user)
        default:
            return state
    }
}

function stats(state = {
                   follow: 0,
                   follower: 0,
                   tweets: 0
               }, action) {
    switch (action.type) {
        case RECEIVE_STATS:
            return Object.assign({}, state, action.stats)
        default:
            return state
    }
}

function auth(state = {
                  isFetching: false,
                  isAuthenticated: !!localStorage.getItem('id_token')
              }, action) {
    switch (action.type) {
        case LOGIN_REQUEST:
            return Object.assign({}, state, {
                isFetching: true,
                isAuthenticated: false,
                user: action.creds
            })
        case LOGIN_SUCCESS:
            return Object.assign({}, state, {
                isFetching: false,
                isAuthenticated: true,
                errorMessage: ''
            })
        case LOGIN_FAILURE:
            return Object.assign({}, state, {
                isFetching: false,
                isAuthenticated: false,
                errorMessage: action.message
            })
        case LOGOUT_SUCCESS:
            return Object.assign({}, state, {
                isFetching: true,
                isAuthenticated: false
            })
        default:
            return state
    }
}

const twitterCloneApp = combineReducers({
    tweets,
    user,
    stats,
    auth
})

export default twitterCloneApp