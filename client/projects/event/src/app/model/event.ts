export class Event {
  private _description: string;
  private _endDate: string;
  private _host: string;
  private _id: string;
  private _location: string;
  private _startDate: string;
  private _title: string;

  constructor(data: any) {
    this._description = data.description;
    this._title = data.title;
    this._endDate = data.endDate;
    this._host = data.host;
    this._id = data.id;
    this._startDate = data.startDate;
    this._location = data.location;
  }

  get description(): string {
    return this._description;
  }

  get endDate(): string {
    return this._endDate;
  }

  get host(): string {
    return this._host;
  }

  get id(): string {
    return this._id;
  }

  get location(): string {
    return this._location;
  }

  get startDate(): string {
    return this._startDate;
  }

  get title(): string {
    return this._title;
  }

  set description(value: string) {
    this._description = value;
  }

  set endDate(value: string) {
    this._endDate = value;
  }

  set host(value: string) {
    this._host = value;
  }

  set id(value: string) {
    this._id = value;
  }

  set location(value: string) {
    this._location = value;
  }

  set startDate(value: string) {
    this._startDate = value;
  }

  set title(value: string) {
    this._title = value;
  }
}
