$(document).ready(function() {
    var seasonTable = $("#seasons").DataTable({
        "ordering": false,
        "info":     false,
        "lengthChange": false,
        "pageLength": 5,
        "sDom": '<"top"i>rt<"bottom"lp><"clear">',
    });
    var timeslotTable = $("#timeslots").DataTable({
        "ordering": false,
        "info":     false,
        "lengthChange": false,
        "sDom": '<"top"i>rt<"bottom"lp><"clear">',
        "pageLength": 5,
        "columnDefs": [
        {
            "targets": [ 0 ],
            "visible": false,
            "searchable": true
        },
        {
            "targets": [ 1,2 ],
            "searchable": false
        },

        ]
    });
    $(".timeslot-filter-season").on( "click", function () {
        $(".timeslotFilterSeason").removeClass("active");
        $(this).addClass("active");
        $("#timeslot-latest").hide();
        $("#timeslot-season-title").text("Season " + $(this).find("td:first-of-type").html())
        timeslotTable
            .columns( 0 )
            .search( $(this).data("seasonid") )
            .draw();

        $("#timeslot-list").show();

        $("#timeslot-list tr").on("click", function () {
            window.location.href = $(this).find("a").attr("href");
        });
    } );
} );
