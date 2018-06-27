// -*- JavaScript -*-
import { QRCode } from 'react-qr-svg';

class Invoice extends React.Component {
  constructor(props) {
    super(props);
    this.state = { invoice: {} };
  }

  componentDidMount() {
    this.serverRequest =
      axios
        .get("/invoice")
        .then((result) => {
           this.setState({ invoice: result.data });
        });
  }

  render() {
      return (
        <div id="qr">{this.state.invoice.payment_request}</div>
      );
  };
}

ReactDOM.render( <Invoice/>, document.querySelector("#root"));
