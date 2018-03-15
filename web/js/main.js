$(document).ready(function () {

    $("#response").html("<strong>Poking server...</strong>");
    $.ajax({
        url: "http://localhost:1234/",
        success: [function (result) {
            $("#response").html("<strong>" + result + "</strong>");
            refreshTransactions();
            refreshInventory();
            refreshLedger();
        }],
        error: [function (result) {
            $("#response").html("<strong>Server is not working, see browser log</strong>");
            console.log(result);
        }]
    });


    function refreshInventory() {
        $.ajax({
            url: "http://localhost:1234/inventory",
            dataType: "json",
            success: [function (result) {
                $('#table-inventory').find('tbody').html("");
                jQuery.each(result.Inventory, function () {


                    var amount = 0;
                    for (var i = 0; i < this.Stacks.length; i++) {
                        amount += this.Stacks[i].Amount;
                    }

                    var newRowContent = "<tr>\n" +
                        "                    <td scope=\"row\">" + this.TypeName + "</td>\n" +
                        "                    <td align='right'>"+amount.toLocaleString()+"</td>\n" +
                        "                </tr>";

                    $('#table-inventory').find('tbody').append(newRowContent);
                });
            }],
            error: [function (result) {
                console.log(result);
            }]
        });
    }

    function refreshLedger() {
        $.ajax({
            url: "http://localhost:1234/ledger",
            dataType: "json",
            success: [function (result) {
                $('#table-ledger').find('tbody').html("");
                jQuery.each(result.Ledger, function () {

                    var newRowContent = "<tr>\n" +
                        "                    <td scope=\"row\">" + this.PlayerName + "</td>\n" +
                        "                    <td align='right'>"+this.Amount.toLocaleString()+" ISK</td>\n" +
                        "                </tr>";

                    $('#table-ledger').find('tbody').append(newRowContent);
                });
            }],
            error: [function (result) {
                console.log(result);
            }]
        });
    }


    function refreshTransactions() {
        $.ajax({
            url: "http://localhost:1234/transactions",
            dataType: "json",
            success: [function (result) {
                $('#uncommitted-transactions').find('tbody').html("");
                jQuery.each(result, function () {

                    var newRowContent = "<tr>\n" +
                        "                    <td scope=\"row\">" + this.CreationDate + "</td>\n" +
                        "                    <td>" + this.PlayerName + "</td>\n" +
                        "                    <td>" + this.Action + "</td>\n" +
                        "                    <td>" + this.TypeName + "</td>\n" +
                        "                    <td align='right'>" + this.Quantity.toLocaleString() + "</td>\n" +
                        "                    <td><input data-id=" + this.Id + " " + (this.MarkedForCorp ? " CHECKED " : "") + " type=\"checkbox\"></td>\n" +
                        "                </tr>";

                    // TODO Add Mark by corp click listener to toggle

                    $('#uncommitted-transactions').find('tbody').append(newRowContent);
                });
            }],
            error: [function (result) {
                console.log(result);
            }]
        });
    }


    $("#submit-parselog").click(function () {
        $.ajax({
            url: "http://localhost:1234/parse",
            method: "POST",
            contentType: "text/plain",
            data: $("#parselog").val(),
            success: [function (result) {
                var responseDiv = $("#responseDiv");
                responseDiv.html("<strong>Parsed transactions!</strong>");
                $("#parselog").val("");
                responseDiv.removeClass("alert-danger");
                responseDiv.removeClass("alert-info");
                responseDiv.addClass("alert-success");
                refreshTransactions();
            }],
            error: [function (result) {
                var responseDiv = $("#responseDiv");
                responseDiv.html("<strong>Could not parse transactions</strong>");
                responseDiv.removeClass("alert-info");
                responseDiv.removeClass("alert-success");
                responseDiv.addClass("alert-danger");
                console.log(result);
            }]
        });
    });
});