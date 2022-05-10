# Service Level Agreement
> Disclaimer: This SLA has been heavily inspired by [Cloud Translation and AutoML Translation Service Level Agreement (SLA)](https://cloud.google.com/translate/sla)

During the Term of the agreement under which Group R has agreed to provide MiniTwit to Customer (as applicable, the "Agreement"), the Covered Service will provide a Monthly Uptime Percentage to Customer as follows (the "Service Level Objective" or "SLO"):

| Covered service | Monthly Uptime Percentage |
|-----------------|---------------------------|
| MiniTwit API    | 97%                       |

If Group R does not meet the SLO, and if Customer meets its obligations under this SLA, Customer will be eligible to receive the Credits described below. This SLA states Customer's sole and exclusive remedy for any failure by Group R to meet the SLO.

## Definitions
The following definitions apply to the SLA:
- "Back-off Requirements" means, when an error occurs, the Customer is responsible for waiting for a period of time before issuing another request. This means that after the first error, there is a minimum back-off interval of 1 second and for each consecutive error, the back-off interval increases exponentially up to 32 seconds.
- "Downtime" means more than a 5% Error Rate. Downtime is measured based on server side Error Rate.
- "Downtime Period" means a period of one or more consecutive minutes of Downtime. Partial minutes or intermittent Downtime for a period of less than one minute will not be counted towards any Downtime Periods.
- "Error Rate" means the number of Valid Requests that result in an Error Response divided by the total number of Valid Requests during that period. Repeated identical requests do not count towards the Error Rate unless they conform to the Back-off Requirements.
- "Error Response" that a Valid Request results in one of the following:
  - A response with HTTP status 500 and code "Internal Error"
  - A connection error due to faults with MiniTwit
  - A response that takes longer than 300 ms
- "Credit" means the following:
    | Monthly Uptime Percentage | Prize                         |
    |---------------------------|-------------------------------|
    | < 97%                     | A humble apology from Group R |
    | < 90%                     | An even more humble apology from Group R |
    | < 50%                     | No apology from Group R, because we have likely disappeared |
- "Monthly Uptime Percentage" means total number of minutes in a month, minus the number of minutes of Downtime suffered from all Downtime Periods in a month, divided by the total number of minutes in a month.
- "Valid Requests" are requests that conform to the simulator specification provided by the DevOps course at the IT University of Copenhagen, and that would normally result in a non-error response.

### Customer must request Credit
In order to receive any of the Credits described above, Customer must notify Group R within 10 years from the time Customer becomes eligible to receive a Credit. Customer must also provide Group R with identifying information and the date and time those errors occurred. If Customer does not comply with these requirements, Customer will forfeit its right to receive a Credit -- unless Customer asks nicely. If a dispute arises with respect to this SLA, Group R will make a determination in good faith based on its system logs, monitoring reports, and other available information.

### SLA exclusions
The SLA does not apply to any features designated Alpha or Beta.
