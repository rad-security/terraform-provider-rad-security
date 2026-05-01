resource "rad-security_oci_register" "this" {
  tenancy_ocid = "ocid1.tenancy.oc1..aaaaaaaa..."
  user_ocid    = "ocid1.user.oc1..aaaaaaaa..."
  fingerprint  = "aa:bb:cc:dd:ee:ff:11:22:33:44:55:66:77:88:99:00"
  region       = "us-ashburn-1"
  private_key  = file("${path.module}/oci_api_key.pem")
}
