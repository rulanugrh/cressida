import { check } from "k6";
import http from "k6/http";
import { Trend } from "k6/metrics";
import { logWaitingTime, URL, logger } from "../util/config";
import { Options } from "k6/options"
import vu from "k6/execution"

type Transporter = {
    vehicle_name: string
    weight: number
    driver_name: string
    distance: string
    price: number
}

const metric = {
    getTransporterByID: new Trend("get_transporter_by_id", true)
}

type Response<T> = {
    code: number,
    msg: string,
    data: T
}

const getTransporterByID = (id: number): Response<Transporter> => {
    const url = URL
    const params = {
        headers: {}
    }

    const res = http.get(`${url}/api/transporters/find/${id}`, params)
    const response = res.json() as { data: Transporter, msg: string }

    logWaitingTime({
        metric: metric.getTransporterByID,
        response: res,
        messageType: "Get Transporter By ID"
    })

    check(res, {
        "200 response code": (r) => r.status === 200,
        "Valid Vehicle Driver Name": (_) => response.data.driver_name !== ''
    })

    return {
        code: res.status,
        msg: response.msg,
        data: response.data
    }
}

export const options: Options = {
    vus: 50,
    iterations: 10,
    duration: '2m'
}

export default () => {

    const randomID = Math.floor(Math.random() * 20)
    logger.info(`Running iteration ${vu.vu.idInInstance} for transporter id: ${randomID}`)
    getTransporterByID(randomID)
}