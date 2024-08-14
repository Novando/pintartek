interface debugParamsType {
  documentURL: string,
  request: {
    postData?: string
  }
}
interface routeType {
  path: string,
  name?: name,
  element: JSX.Element,
  needAuth?: bool,
}
interface responseType {
  data: any,
  message: string,
  count?: number,
}
interface FilterQueryType {
  page?: number,
  size?: number,
}
interface IwkbuSearchType {
  period?: string, // ISO Format
  official?: string,
}
interface BaseIwkbuDataType {
  note:	string,
  reasonId: number,
}
interface IwkbuDataType extends BaseIwkbuDataType {
  id: number,
  nopol: string,
  lastPrice: number,
  lastRecorded: string, // ISO Format
  lastIwkbu: string, // ISO Format
  lastSwdkllj: string, // ISO Format
  currPrice: number,
  currRecorded: string, // ISO Format
  currIwkbu: string, // ISO Format
  currSwdkllj: string, // ISO Format
  diffPrice: number,
}
interface IwkbuYearUpdateDataType extends BaseIwkbuDataType {
  price: number,
  recorded: string, // ISO Format
  inputedOn: string, // ISO Format
  iwkbu: string, // ISO Format
  swdkllj: string, // ISO Format
}
interface IwkbuYearCreateDataType extends BaseIwkbuDataType {
  price: number,
  nopol: string,
  kantor: string,
  recorded: string, // ISO Format
  inputedOn: string, // ISO Format
  iwkbu: string, // ISO Format
  swdkllj: string, // ISO Format
}
interface OfficeQueryType {
  search: string,
  hq?: string,
}
interface OfficialDataType {
  id: number,
  name: string,
}
interface HqDataType {
  id: number,
  code: string,
}
interface ConversionReportRequestDataType {
  hq?: string,
  till?: string, // ISO Format
  since?: string, // ISO Format
  office?: string,
}
interface ReasonDataType {
  id: number,
  score: number,
  name: string,
  desc: string,
}
interface ReportConvDetailDataType {
  name: string,
  price: number,
  qty: number,
  score: number,
}
interface ReportOfficialDetailDataType {
  hq: string,
  office: string,
  currIncome: number,
  currQty: number,
  lastIncome: number,
  lastQty: number,
  sumIncome: number,
  sumQty: number,
}
interface ReportConvSummaryDataType {
  currIncome: number,
  currNopol: number,
  lastIncome: number,
  lastNopol: number,
}
interface ReportConvDataType {
  details: ReportConvDetailDataType[],
  summary: ReportConvSummaryDataType,
}