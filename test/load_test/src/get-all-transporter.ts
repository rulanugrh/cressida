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
    getAllTransporter: new Trend("get_all_transporter", true)
}

type Response<T> = {
    code: number,
    msg: string,
    data: T[]
}

const getAllVehicle = (): Response<Transporter> => {
    const url = URL

    const res = http.get(`${url}/api/transporters/get`)
    const response = res.json() as { data: Transporter[], msg: string }

    logWaitingTime({
        metric: metric.getAllTransporter,
        response: res,
        messageType: "Get All Transporter"
    })


    check(res, {
        "200 response code": (r) => r.status === 200,
        "Valid message response": (_) => response.msg === "success found transporter"
    })

    return {
        code: res.status,
        msg: response.msg,
        data: response.data
    }
}

export const options: Options = {
    vus: 10,
    iterations: 5,
    duration: '2m'
}

export default () => {
    logger.info(`Running iteration ${vu.vu.idInInstance}`)
    getAllVehicle()
}