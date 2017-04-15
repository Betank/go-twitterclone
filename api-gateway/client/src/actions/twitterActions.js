import fetch from 'isomorphic-fetch'

export const ADD_TWEET = 'ADD_TWEET'
export const ADD_TWEET_FAILED = 'ADD_TWEET_FAILED'
export const ADD_TWEET_SUCCESS = 'ADD_TWEET_SUCCESS'
export const RECEIVE_TWEETS = 'RECEIVE_TWEETS'
export const RECEIVE_PROFILE = 'RECEIVE_PROFILE'
export const RECEIVE_STATS = 'RECEIVE_STATS'

export const LOGIN_REQUEST = 'LOGIN_REQUEST'
export const LOGIN_SUCCESS = 'LOGIN_SUCCESS'
export const LOGIN_FAILURE = 'LOGIN_FAILURE'

export const LOGOUT_REQUEST = 'LOGOUT_REQUEST'
export const LOGOUT_SUCCESS = 'LOGOUT_SUCCESS'

export const addTweet = (text) => {
    let token = localStorage.getItem('id_token')
    if (!token) {
        throw "unable to find token"
    }

    return dispatch => {
        dispatch({ type: ADD_TWEET })
        return fetch('/api/tweet/', {
            method: 'POST',
            headers: {
                'Content-Type' : 'application/text; charset=UTF-8',
                Authorization: `Bearer ${token}`
            },
            body: text
        })
        .then(res => {
            if (res.status != 200) {
                dispatch({ type: ADD_TWEET_FAILED })
            } else {
                dispatch({ type: ADD_TWEET_SUCCESS })
            }
        })
        .catch(err => {
            dispatch({ type: ADD_TWEET_FAILED })
        })
    }
}

export function fetchTweets() {
    let token = localStorage.getItem('id_token')
    if (!token) {
        throw "unable to find token"
    }


    return dispatch => {
        return fetch('/api/tweets/user/', { headers: {Authorization: `Bearer ${token}`}})
        .then(response => response.json())
        .then(json => dispatch(receiveTweets(json)))
    }
}

function receiveTweets(tweets) {
        return {
        type: RECEIVE_TWEETS,
        tweets: tweets
    }
}

export function fetchProfile() {
    let token = localStorage.getItem('id_token')
    if (!token) {
        throw "unable to find token"
    }

    return dispatch => {
        return fetch('/api/user/', { headers: {Authorization: `Bearer ${token}`}})
        .then(response => response.json())
        .then(json => dispatch(receiveProfile(json)))
    }
}

function receiveProfile(user) {
        return {
        type: RECEIVE_PROFILE,
        user: user
    }
}

export function fetchStatistics() {
    let token = localStorage.getItem('id_token')
    if (!token) {
        throw "unable to find token"
    }

    return dispatch => {
        return fetch('/api/stats/', { headers: {Authorization: `Bearer ${token}`}})
        .then(response => response.json())
        .then(json => dispatch(receiveStatistics(json)))
    }
}

function receiveStatistics(stats) {
    return {
        type: RECEIVE_STATS,
        stats: stats
    }
}

function requestLogin(creds) {
    return {
        type: LOGIN_REQUEST,
        isFetching: true,
        isAuthenticated: false,
        creds
    }
}

function receiveLogin(user) {
    return {
        type: LOGIN_SUCCESS,
        isFetching: false,
        isAuthenticated: true,
        id_token: user.token
    }
}

function loginError(message) {
    return {
        type: LOGIN_FAILURE,
        isFetching: false,
        isAuthenticated: false,
        message
    }
}

export function loginUser(creds) {
    let config = {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded'},
        body: 'username=${creds.username}&password=${creds.password}'
    }

    return dispatch => {
        dispatch(requestLogin(creds))

        return fetch('/api/login/', config)
            .then(response => response.json().then(user => ({user, response})))
            .then(({user, response}) => {
                if (!response.ok) {
                    dispatch(loginError(user.message))
                    return Promise.reject(user)
                } else {
                    localStorage.setItem('id_token', user.token)
                    dispatch(receiveLogin(user))
                }
            })
            .catch(err => console.log('Error: ', err))
    }
}

function requestLogout() {
    return {
        type: LOGOUT_REQUEST,
        isFetching: true,
        isAuthenticated: true
    }
}

function receiveLogout() {
    return {
        type: LOGOUT_SUCCESS,
        isFetching: false,
        isAuthenticated: false
    }
}

export function logoutUser() {
    return dispatch => {
        dispatch(requestLogout())
        localStorage.removeItem('id_token')
        dispatch(receiveLogout())
    }
}