package helpers

import (
	"fmt"
)

func SendOTP(email, otpType string, username string) error {
	otp := GenerateOTP()

	fmt.Println("Generated OTP:", otp)
	fmt.Printf("username: %s\n", username)

	var mailSubject string
	if otpType == "forgot" {
		mailSubject = "Password Reset Request"
	} else {
		mailSubject = "Complete your registration"
	}

	fmt.Printf("Sending %s OTP to %s\n", otp, email)

	content := fmt.Sprintf(`
	<!DOCTYPE html>
	<html lang="en" xmlns:o="urn:schemas-microsoft-com:office:office" xmlns:v="urn:schemas-microsoft-com:vml">
	<head>
		<title></title>
		<meta content="text/html; charset=utf-8" http-equiv="Content-Type"/>
		<meta content="width=device-width, initial-scale=1.0" name="viewport"/>
		<style>
			* {
				box-sizing: border-box;
			}

			body {
				margin: 0;
				padding: 0;
			}

			a[x-apple-data-detectors] {
				color: inherit !important;
				text-decoration: inherit !important;
			}

			#MessageViewBody a {
				color: inherit;
				text-decoration: none;
			}

			p {
				line-height: inherit
			}

			.desktop_hide,
			.desktop_hide table {
				mso-hide: all;
				display: none;
				max-height: 0px;
				overflow: hidden;
			}

			.image_block img+div {
				display: none;
			}

			@media (max-width:620px) {
				.desktop_hide table.icons-inner {
					display: inline-block !important;
				}

				.icons-inner {
					text-align: center;
				}

				.icons-inner td {
					margin: 0 auto;
				}

				.row-content {
					width: 100% !important;
				}

				.mobile_hide {
					display: none;
				}

				.stack .column {
					width: 100%;
					display: block;
				}

				.mobile_hide {
					min-height: 0;
					max-height: 0;
					max-width: 0;
					overflow: hidden;
					font-size: 0px;
				}

				.desktop_hide,
				.desktop_hide table {
					display: table !important;
					max-height: none !important;
				}

				.row-1 .column-1 .block-6.text_block td.pad {
					padding: 30px 30px 25px !important;
				}

				.row-1 .column-1 .block-5.heading_block h1 {
					font-size: 39px !important;
				}

				.row-1 .column-1 .block-4.text_block td.pad {
					padding: 30px 20px 20px !important;
				}

				.row-1 .column-1 .block-3.text_block td.pad {
					padding: 10px 10px 15px 15px !important;
				}

				.row-1 .column-1 .block-2.image_block td.pad {
					padding: 30px 30px 30px 20px !important;
				}
			}
		</style>
	</head>
	<body style="background-color: #76c8f7; margin: 0; padding: 0; -webkit-text-size-adjust: none; text-size-adjust: none;">
		<table border="0" cellpadding="0" cellspacing="0" class="nl-container" role="presentation" style="mso-table-lspace: 0pt; mso-table-rspace: 0pt; background-color: #76c8f7;" width="100%">
			<tbody>
				<tr>
					<td>
						<table align="center" border="0" cellpadding="0" cellspacing="0" class="row row-1" role="presentation" style="mso-table-lspace: 0pt; mso-table-rspace: 0pt; background-color: #f5f5f5; background-image: url(''); background-position: center top; background-repeat: repeat;" width="100%">
							<tbody>
								<tr>
									<td>
										<table align="center" border="0" cellpadding="0" cellspacing="0" class="row-content stack" role="presentation" style="mso-table-lspace: 0pt; mso-table-rspace: 0pt; background-color: #ffffff; color: #000000; width: 600px;" width="600">
											<tbody>
												<tr>
													<td class="column column-1" style="mso-table-lspace: 0pt; mso-table-rspace: 0pt; font-weight: 400; text-align: left; padding-bottom: 15px; padding-left: 10px; padding-right: 10px; padding-top: 5px; vertical-align: top; border-top: 0px; border-right: 0px; border-bottom: 0px; border-left: 0px;" width="100%">
														<div class="spacer_block block-1" style="height:39px;line-height:39px;font-size:1px;"> </div>
														<table border="0" cellpadding="0" cellspacing="0" class="text_block block-3" role="presentation" style="mso-table-lspace: 0pt; mso-table-rspace: 0pt; word-break: break-word;" width="100%">
															<tr>
																<td class="pad" style="padding-bottom:15px;padding-left:40px;padding-right:10px;padding-top:10px;">
																	<div style="font-family: Tahoma, Verdana, sans-serif">
																		<div style="font-size: 14px; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; mso-line-height-alt: 28px; color: #000000; line-height: 2;">
																			<p style="margin: 0; font-size: 15px; text-align: left; mso-line-height-alt: 30px;"><span style="font-size:15px;">Hi there, %s</span></p>
																		</div>
																	</div>
																</td>
															</tr>
														</table>
														<table border="0" cellpadding="20" cellspacing="0" class="heading_block block-5" role="presentation" style="mso-table-lspace: 0pt; mso-table-rspace: 0pt;" width="100%">
															<tr>
																<td class="pad">
																	<h1 style="margin: 0; color: #5327b5; direction: rtl; font-family: 'Varela Round', 'Trebuchet MS', Helvetica, sans-serif; font-size: 53px; font-weight: 700; letter-spacing: 24px; line-height: 120%; text-align: center; margin-top: 0; margin-bottom: 0;"><strong>%s</strong></h1>
																</td>
															</tr>
														</table>
														<table border="0" cellpadding="45" cellspacing="0" class="text_block block-6" role="presentation" style="mso-table-lspace: 0pt; mso-table-rspace: 0pt; word-break: break-word;" width="100%">
															<tr>
																<td class="pad">
																	<div style="font-family: Arial, sans-serif">
																		<div class="" style="font-size: 14px; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; mso-line-height-alt: 28px; color: #000000; line-height: 2;">
																			<p style="margin: 0; font-size: 15px; mso-line-height-alt: 30px;"><span style="font-size:15px;">If you did not associate your address with your BlackStar Enterprise account, please ignore this message and do not click on the link above.</span></p>
																			<p style="margin: 0; font-size: 15px; mso-line-height-alt: 30px;"><span style="font-size:15px;">If you experience any issues, don't hesitate to reach out to our support team 👉 <a href="mailto:supporttaraqinvestment@gmail.com" rel="noopener" style="text-decoration: underline; color: #ef0c0c;" target="_blank" title="hello@use.com">here</a></span></p>
																			<p style="margin: 0; font-size: 15px; mso-line-height-alt: 30px;"><span style="font-size:15px;">Best Regards</span></p>
																			<p style="margin: 0; font-size: 15px; mso-line-height-alt: 30px;"><span style="font-size:15px;color:#5327b5;"><strong><span style="">Team Taraq Investment App.</span></strong></span></p>
																		</div>
																	</div>
																</td>
															</tr>
														</table>
													</td>
												</tr>
											</tbody>
										</table>
									</td>
								</tr>
							</tbody>
						</table>
					</td>
				</tr>
			</tbody>
		</table>
	</body>
	</html>
	`, username, otp)

	err := sendMail(email, mailSubject, content)
	if err != nil {
		return err
	}

	StoreOTP(email, otp)
	return nil
}