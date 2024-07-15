import { check } from "k6";
import http from "k6/http";
import { Trend } from "k6/metrics";
import { logWaitingTime, URL, logger } from "../util/config";
import { Options } from "k6/options"
import vu from "k6/execution"

type Vehicle = {
    id: number
    name: string
    description: string
}

const metric = {
    getVehicleByID: new Trend("get_vehicle_by_id", true)
}

type Response<T> = {
    code: number,
    msg: string,
    data: T
}

const getVehicleByID = (id: number): Response<Vehicle> => {
    const url = URL
    const params = {
        headers: {}
    }

    const res = http.get(`${url}/api/vehicles/find/${id}`, params)
    const response = res.json() as { data: Vehicle, msg: string }

    logWaitingTime({
        metric: metric.getVehicleByID,
        response: res,
        messageType: "Get Vehicle By ID"
    })

    check(res, {
        "200 response code": (r) => r.status === 200,
        "Valid Vehicle ID": (_) => response.data.id === id
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
    logger.info(`Running iteration ${vu.vu.idInInstance} for vehicle id: ${randomID}`)
    getVehicleByID(randomID)
}