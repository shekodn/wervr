var App = React.createClass({
  componentWillMount: function() {
    this.setupAjax();
    this.parseHash();
    this.setState();
  },
  // Add access_token if available with each XHR request to API
  setupAjax: function() {
    $.ajaxSetup({
      'beforeSend': function(xhr) {
        if (localStorage.getItem('access_token')) {
          xhr.setRequestHeader('Authorization',
                'Bearer ' + localStorage.getItem('access_token'));
        }
      }
    });
  },
  // Extract the access_token and id_token from Auth0 Callback after login
  parseHash: function(){
    this.auth0 = new auth0.WebAuth({
      domain:       AUTH0_DOMAIN,
      clientID:     AUTH0_CLIENT_ID
    });
    this.auth0.parseHash(window.location.hash, function(err, authResult) {
      if (err) {
        return console.log(err);
      }
      if(authResult !== null && authResult.accessToken !== null && authResult.idToken !== null){
        localStorage.setItem('access_token', authResult.accessToken);
        localStorage.setItem('id_token', authResult.idToken);
        localStorage.setItem('profile', JSON.stringify(authResult.idTokenPayload));
    window.location = window.location.href.substr(0, window.location.href.indexOf('#'))
      }
    });
  },
  // Set user login state
  setState: function(){
    var idToken = localStorage.getItem('id_token');
    if(idToken){
      this.loggedIn = true;
    } else {
      this.loggedIn = false;
    }
  },
  render: function() {

    if (this.loggedIn) {
      return (<LoggedIn />);
    } else {
      return (<Home />);
    }
  }
});



var LoggedIn = React.createClass({
  // If a user logs out we will remove their tokens and profile info
  logout : function(){
    localStorage.removeItem('id_token');
    localStorage.removeItem('access_token');
    localStorage.removeItem('profile');
    location.reload();
  },
  getInitialState: function() {
    return {
      products: []
    }
  },
  // Once this components mounts, we will make a call to the API to get the product data
  componentDidMount: function() {
    this.serverRequest = $.get('http://localhost:3000/products', function (result) {
      this.setState({
        products: result,
      });
    }.bind(this));
  },

  render: function() {
    return (
      <div className="col-lg-12">
        <span className="pull-right"><a onClick={this.logout}>Log out</a></span>
        <h2>Welcome to We R VR</h2>
        <p>Below you'll find the latest games that need feedback. Please provide honest feedback so developers can make the best games.</p>
        <div className="row">

        {this.state.products.map(function(product, i){
          return <Product key={i} product={product} />
        })}

        </div>
      </div>);
  }
});



var Home = React.createClass({
  // Clicking the login link will redirect the user to the Hosted Lock page
  // where the user will provide their credentials and will then be redirected
  // back to the application
  authenticate: function(){
    this.webAuth = new auth0.WebAuth({
      domain:       AUTH0_DOMAIN,
      clientID:     AUTH0_CLIENT_ID,
      scope:        'openid profile',
      audience:     AUTH0_API_AUDIENCE,
      responseType: 'token id_token',
      redirectUri : AUTH0_CALLBACK_URL
    });
    this.webAuth.authorize();
  },
  render: function() {
    return (
    <div className="container">
      <div className="col-xs-12 jumbotron text-center">
        <h1>We R VR</h1>
        <p>Provide valuable feedback to VR experience developers.</p>
        <a onClick={this.authenticate} className="btn btn-primary btn-lg btn-login btn-block">Sign In</a>
      </div>
    </div>);
  }
});


var Product = React.createClass({
  // The upvote and downvote functions will send an API request to the backend
  // These calls too will send the access_token and be verified before a success response
  // is returned
  upvote : function(){
    var product = this.props.product;
    this.serverRequest = $.post('http://localhost:3000/products/' + product.Slug + '/feedback', {vote : 1}, function (result) {
      this.setState({voted: "Upvoted"})
    }.bind(this));
  },
  downvote: function(){
    var product = this.props.product;
    this.serverRequest = $.post('http://localhost:3000/products/' + product.Slug + '/feedback', {vote : -1}, function (result) {
      this.setState({voted: "Downvoted"})
    }.bind(this));
  },
  getInitialState: function() {
    return {
      voted: null
    }
  },
  render : function(){
    return(
    <div className="col-xs-4">
      <div className="panel panel-default">
        <div className="panel-heading">{this.props.product.Name} <span className="pull-right">{this.state.voted}</span></div>
        <div className="panel-body">
          {this.props.product.Description}
        </div>
        <div className="panel-footer">
          <a onClick={this.upvote} className="btn btn-default">
            <span className="glyphicon glyphicon-thumbs-up"></span>
          </a>
          <a onClick={this.downvote} className="btn btn-default pull-right">
            <span className="glyphicon glyphicon-thumbs-down"></span>
          </a>
        </div>
      </div>
    </div>);
  }
})


ReactDOM.render(<App />, document.getElementById('app'));
