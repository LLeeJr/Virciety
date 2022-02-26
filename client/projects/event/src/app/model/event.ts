export class Event {
  private _description: string;
  private _endDate: string;
  private readonly _host: string;
  private readonly _id: string;
  private _location: string;
  private _startDate: string;
  private _title: string;
  private _subscribers: string[];
  private _startTime: string | null;
  private _endTime: string | null;
  private _attendees: string[];
  private _currentlyAttended: boolean;

  constructor(data: any) {
    this._description = data.description;
    this._title = data.title;
    this._host = data.host;
    this._id = data.id;
    this._location = data.location;
    this._startDate = data.startDate;
    this._endDate = data.endDate;
    this._subscribers = data.subscribers;
    this._attendees = data.attendees;
    this._currentlyAttended = data.currentlyAttended;

    if (data.startDate.endsWith('M') && data.endDate.endsWith('M')) {
      this._startTime = data.startDate.split(',')[1].trim();
      this._endTime = data.endDate.split(',')[1].trim();
    } else {
      this._startTime = null;
      this._endTime = null;
    }
  }

  get currentlyAttended(): boolean {
    return this._currentlyAttended;
  }

  get attendees(): string[] {
    return this._attendees;
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

  get startTime(): string | null {
    return this._startTime;
  }

  get endTime(): string | null {
    return this._endTime;
  }

  get subscribers(): string[] {
    return this._subscribers;
  }

  set currentlyAttended(value: boolean) {
    this._currentlyAttended = value;
  }

  set attendees(value: string[]) {
    this._attendees = value;
  }

  set description(value: string) {
    this._description = value;
  }

  set subscribers(value: string[]) {
    this._subscribers = value;
  }

  set endDate(value: string) {
    this._endDate = value;
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

  set startTime(value: string | null) {
    this._startTime = value;
  }

  set endTime(value: string | null) {
    this._endTime = value;
  }
}
