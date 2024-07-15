import { check } from "k6";
import { Trend } from "k6/metrics";
import { logWaitingTime, URL, logger } from "../util/config";
import { Options } from "k6/options"
import vu from "k6/execution"
import http from "k6/http";
import { faker } from "@faker-js/faker";

type Transporter = {
    vehicle_name: string
    weight: number
    driver_name: string
    distance: string
    price: number
}

type Response<T> = {
    code: number
    msg: string
    data: T
}

const metric = {
    createTransporter: new Trend("create_request_transporter", true),
    userLogin: new Trend("user_login", true)

}

const loginRequest = (): string => {
    const url = URL

    const response = http.post(`${url}/api/users/login`, {
        email: process.env.ADMIN_EMAIL,
        password: process.env.ADMIN_PASSWORD
    })

    const jsonRresponse = response.json() as { data: string, msg: string }

    logWaitingTime({
        metric: metric.userLogin,
        response: response,
        messageType: "User Login"
    })

    check(response, {
        'Success Login': (r) => r.status === 201,
        'Message Valid': (_) => jsonRresponse.msg === "success login user",
        'loggin success': () => jsonRresponse.data !== ''
    })

    return jsonRresponse.data
}

const createTransporter = (token: string, driver_id: number, vehicle_type: number, max_weight: number, max_distance: number, price: number): Response<Transporter> => {
    const url = URL
    const headers = {
        'Content-Type': 'application/json',
        'Authorization': `${token}`
    }

    const payload = JSON.stringify({
        driver_id: driver_id,
        vehicle_type: vehicle_type,
        max_weight: max_weight,
        max_distance: max_distance,
        price: price
    })

    const res = http.post(`${url}/api/vehicles/add`, payload, { headers: headers})
    const response = res.json() as { data: Transporter, msg: string }

    logWaitingTime({
        metric: metric.createTransporter,
        response: response,
        messageType: "Create Vehicle"
    })

    check(res, {
        'Success Post Data': (r) => r.status === 201,
        'Message Valid': (_) => response.msg === "success add vehicle",
        'Data Valid': (_) => response.data.driver_name !== ''
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
    duration: '2s'
}

export default () => {
    const token = loginRequest()
    const res = createTransporter(token, faker.number.int({ min: 1, max: 10}), faker.number.int({ min: 1, max: 10}), faker.number.int({ min: 20, max: 50}), faker.number.int({ min: 20, max: 50}), faker.number.int({ min: 150000, max: 300000}))
    logger.info(`Running iteration ${vu.vu.idInInstance} for vehicle name: ${res.data.vehicle_name}`)

}