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
    getAllVehicle: new Trend("get_all_vehicle", true)
}

type Response<T> = {
    code: number,
    msg: string,
    data: T[]
}

const getAllVehicle = (): Response<Vehicle> => {
    const url = URL

    const res = http.get(`${url}/api/vehicles/get`)
    const response = res.json() as { data: Vehicle[], msg: string }

    logWaitingTime({
        metric: metric.getAllVehicle,
        response: res,
        messageType: "Get All Vehicle"
    })


    check(res, {
        "200 response code": (r) => r.status === 200,
        "Valid message response": (_) => response.msg === "success found vehicle"
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
    const response = getAllVehicle()
    for (let res of response.data) {
        logger.info(`Running iteration ${vu.vu.idInInstance} for vehicle id: ${res.id}`)
    }
}