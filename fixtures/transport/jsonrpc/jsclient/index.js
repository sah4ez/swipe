// Code generated by Swipe v1.11.4. DO NOT EDIT.

class JSONRPCClient {
  /**
   *
   * @param {string} rgt
   */
  constructor(tgt) {
    this.tgt = tgt;
    this._requestID = 0;
    this._scheduleRequests = {};
    this._commitTimerID = null;
    this._beforeRequest = null;
  }
  beforeRequest(fn) {
    this._beforeRequest = fn;
  }
  __scheduleCommit() {
    if (this._commitTimerID) {
      clearTimeout(this._commitTimerID);
    }
    this._commitTimerID = setTimeout(() => {
      this._commitTimerID = null;
      const scheduleRequests = { ...this._scheduleRequests };
      this._scheduleRequests = {};

      let requests = [];

      for (let key in scheduleRequests) {
        requests.push(scheduleRequests[key].request);
      }

      this.__doRequest(requests)
        .then(responses => {
          for (var i = 0; i < responses.length; i++) {
            if (responses[i].error) {
              scheduleRequests[responses[i].id].reject(responses[i].error);
              continue;
            }
            scheduleRequests[responses[i].id].resolve(responses[i].result);
          }
        })
        .catch(e => {
          for (var key in requests) {
            scheduleRequests[key].reject(e);
          }
        });
    }, 0);
  }
  makeJSONRPCRequest(id, method, params) {
    return {
      jsonrpc: "2.0",
      id: id,
      method: method,
      params: params
    };
  }
  __scheduleRequest(method, params) {
    var p = new Promise((resolve, reject) => {
      const request = this.makeJSONRPCRequest(
        this.__requestIDGenerate(),
        method,
        params
      );
      this._scheduleRequests[request.id] = {
        request,
        resolve,
        reject
      };
    });
    this.__scheduleCommit();
    return p;
  }

  __doRequest(request) {
    let header = {
      "Content-Type": "application/json;charset=utf-8"
    };
    if (this._beforeRequest) {
      let ctx = this._beforeRequest({ header });
      if (ctx) {
        if (ctx.header) {
          header = {
            ...header,
            ...ctx.header
          };
        }
      }
    }
    return fetch(this.tgt, {
      method: "POST",
      headers: header,
      body: JSON.stringify(request)
    }).then(response => {
      return response.json();
    });
  }
  __requestIDGenerate() {
    return ++this._requestID;
  }
}
export default class ClientServiceInterface extends JSONRPCClient {
  /**
   * @param {number} id
   * @param {string} name
   * @param {string} fname
   * @param {number} price
   * @param {number} n
   * @return {PromiseLike<data: {id: string, name: string, password: string, point: {type: string, coordinates: Array.<number>}, lastSeen: string, photo: Array.<number>, profile: {phone: string}}>}
   **/
  get(id, name, fname, price, n) {
    return this.__scheduleRequest("get", {
      id: id,
      name: name,
      fname: fname,
      price: price,
      n: n
    });
  }
  /**
   * @return {PromiseLike<Array.<{id: string, name: string, password: string, point: {type: string, coordinates: Array.<number>}, lastSeen: string, photo: Array.<number>, profile: {phone: string}}>>}
   **/
  getAll() {
    return this.__scheduleRequest("getAll", {});
  }
  /**
   * @param {Object.<string, *>} data
   * @param {*} ss
   * @return {PromiseLike<Object.<string, string>>}
   **/
  testMethod(data, ss) {
    return this.__scheduleRequest("testMethod", { data: data, ss: ss });
  }
  /**
   * @param {string} name
   * @param {Array.<number>} data
   **/
  create(name, data) {
    return this.__scheduleRequest("create", { name: name, data: data });
  }
  /**
   * @param {number} id
   * @return {PromiseLike<a: string, b: string>}
   **/
  delete(id) {
    return this.__scheduleRequest("delete", { id: id });
  }
}